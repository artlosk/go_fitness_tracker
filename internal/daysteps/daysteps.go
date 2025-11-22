// Package daysteps
// Отвечает за учёт активности в течение дня.
// Он собирает переданную информацию в виде строк, парсит их и выводит информацию о количестве шагов, пройденной дистанции и потраченных калориях.
package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// Функция парсит строку формата "678,0h50m"
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return 0, 0, errors.New("неверные входные данные")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("непраильный формат шагов: %w", err)
	}

	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	durationWalk, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}

	if durationWalk <= 0 {
		return 0, 0, errors.New("продолжительность должна быть больше 0")
	}

	return steps, durationWalk, nil
}

// DayActionInfo
// Функция вычисляет дистанцию в километрах и количество потраченных калорий и возвращать строку в таком виде
func DayActionInfo(data string, weight, height float64) string {
	steps, durationWalk, err := parsePackage(data)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	distance := (float64(steps) * stepLength) / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, durationWalk)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distance, calories,
	)

	return result
}
