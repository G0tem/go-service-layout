package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Config struct {
	URL              string        `envconfig:"HEALTHCHECK_URL" default:"http://app:8080/api/v1/healthz"`
	Interval         time.Duration `envconfig:"HEALTHCHECK_INTERVAL" default:"30s"`
	Timeout          time.Duration `envconfig:"HEALTHCHECK_TIMEOUT" default:"5s"`
	FailureThreshold int           `envconfig:"HEALTHCHECK_FAILURE_THRESHOLD" default:"3"`
}

type HealthChecker struct {
	config            Config
	client            *http.Client
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
	isHealthy         bool
	mu                sync.RWMutex
	failureCount      int
	lastCheckTime     time.Time
	lastError         error
	checkCount        int64
	successCount      int64
	failureCountTotal int64
}

type HealthStatus struct {
	IsHealthy     bool      `json:"is_healthy"`
	LastCheckTime time.Time `json:"last_check_time"`
	LastError     error     `json:"last_error,omitempty"`
	CheckCount    int64     `json:"check_count"`
	SuccessCount  int64     `json:"success_count"`
	FailureCount  int64     `json:"failure_count"`
	CurrentStreak int       `json:"current_streak"`
}

func New(cfg Config) *HealthChecker {
	ctx, cancel := context.WithCancel(context.Background())

	hc := &HealthChecker{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
		ctx:       ctx,
		cancel:    cancel,
		isHealthy: true, // По умолчанию считаем здоровым до первой проверки
	}

	return hc
}

func (hc *HealthChecker) Start() {
	hc.wg.Add(1)
	go hc.run()
	log.Info().
		Str("url", hc.config.URL).
		Dur("interval", hc.config.Interval).
		Msg("Health checker started")
}

func (hc *HealthChecker) Stop() {
	log.Info().Msg("Stopping health checker...")
	hc.cancel()
	hc.wg.Wait()
	log.Info().Msg("Health checker stopped")
}

func (hc *HealthChecker) run() {
	defer hc.wg.Done()

	ticker := time.NewTicker(hc.config.Interval)
	defer ticker.Stop()

	// Первая проверка сразу после запуска
	hc.check()

	for {
		select {
		case <-hc.ctx.Done():
			return
		case <-ticker.C:
			hc.check()
		}
	}
}

func (hc *HealthChecker) check() {
	hc.mu.Lock()
	hc.checkCount++
	hc.lastCheckTime = time.Now()
	hc.mu.Unlock()

	start := time.Now()
	resp, err := hc.client.Get(hc.config.URL)

	duration := time.Since(start)

	hc.mu.Lock()
	defer hc.mu.Unlock()

	if err != nil {
		hc.failureCount++
		hc.failureCountTotal++
		hc.lastError = err
		hc.isHealthy = false

		log.Error().
			Err(err).
			Dur("duration", duration).
			Int("failure_count", hc.failureCount).
			Msg("Health check failed")

		// Логгируем только если превышен порог неудач
		if hc.failureCount >= hc.config.FailureThreshold {
			log.Warn().
				Err(err).
				Int("consecutive_failures", hc.failureCount).
				Msg("Health check consecutive failures threshold exceeded")
		}

		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		hc.failureCount++
		hc.failureCountTotal++
		hc.lastError = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		hc.isHealthy = false

		log.Warn().
			Int("status_code", resp.StatusCode).
			Dur("duration", duration).
			Int("failure_count", hc.failureCount).
			Msg("Health check returned non-OK status")

		return
	}

	// Успешная проверка
	hc.successCount++
	hc.failureCount = 0 // Сбрасываем счетчик неудач
	hc.lastError = nil
	hc.isHealthy = true

	log.Debug().
		Dur("duration", duration).
		Msg("Health check passed")
}

func (hc *HealthChecker) IsHealthy() bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	return hc.isHealthy
}

func (hc *HealthChecker) GetStatus() HealthStatus {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	streak := int(hc.successCount - hc.failureCountTotal)
	if streak < 0 {
		streak = 0
	}

	return HealthStatus{
		IsHealthy:     hc.isHealthy,
		LastCheckTime: hc.lastCheckTime,
		LastError:     hc.lastError,
		CheckCount:    hc.checkCount,
		SuccessCount:  hc.successCount,
		FailureCount:  hc.failureCountTotal,
		CurrentStreak: streak,
	}
}

// WaitForHealthy ожидает пока сервис станет здоровым или таймаут
func (hc *HealthChecker) WaitForHealthy(ctx context.Context, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			if hc.IsHealthy() {
				return true
			}
		}
	}
}
