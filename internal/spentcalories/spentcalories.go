// Package spentcalories
// Обрабатывает переданную информацию и рассчитывает потраченные калории в зависимости от вида активности — бега или ходьбы.
// И возвращает информацию обо всех тренировках.
package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// Функция принимает строку с данными формата "3456,Ходьба,3h00m" и парсит данные
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверные входные данные")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
	}

	activityType := parts[1]

	durationWalk, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	if durationWalk <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть больше 0")
	}

	return steps, activityType, durationWalk, nil
}

// Парсит строку, переводит данные из строки в соответствующие типы
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return (float64(steps) * stepLength) / mInKm
}

// Средняя скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

// TrainingInfo
// Отвечает за подсчет и вывод данных
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, durationWalk, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, durationWalk)

	var calories float64
	switch activityType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, durationWalk)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, durationWalk)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activityType)
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType,
		durationWalk.Hours(),
		dist,
		speed,
		calories,
	)

	return result, nil
}

// RunningSpentCalories
// Считает потраченные калории при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateParams(steps, weight, height, duration); err != nil {
		return 0, err
	}
	return (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH, nil
}

// WalkingSpentCalories
// Считает потраченные калории при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateParams(steps, weight, height, duration); err != nil {
		return 0, err
	}
	return ((weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH) * walkingCaloriesCoefficient, nil
}

func validateParams(steps int, weight, height float64, duration time.Duration) error {
	checks := []struct {
		ok      bool
		message string
	}{
		{steps > 0, "количество шагов должно быть больше 0"},
		{weight > 0, "вес должен быть больше 0"},
		{height > 0, "рост должен быть больше 0"},
		{duration > 0, "продолжительность должна быть больше 0"},
	}

	for _, check := range checks {
		if !check.ok {
			return errors.New(check.message)
		}
	}
	return nil
}
