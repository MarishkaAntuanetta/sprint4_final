package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // Количество метров в одном километре.
	minInH                     = 60   // Количество минут в одном часе.
	stepLengthCoefficient      = 0.45 // Коэффициент для расчета длины шага на основе роста пользователя.
	walkingCaloriesCoefficient = 0.5  // Коэффициент для расчета калорий при ходьбе.
)

// parseTraining разбирает входные данные и извлекает количество шагов, тип тренировки и её продолжительность.
func parseTraining(data string) (int, string, time.Duration, error) {
	strData := strings.Split(data, ",")
	if len(strData) != 3 {
		return 0, "", 0, errors.New("invalid data format")
	}

	step, err := strconv.Atoi(strData[0])
	if err != nil {
		return 0, "", 0, err
	}
	if step <= 0 {
		return 0, "", 0, errors.New("invalid step value")
	}
	duration, err := time.ParseDuration(strData[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("invalid activity duration")
	}
	return step, strData[1], duration, nil
}

// distance вычисляет дистанцию, пройденную пользователем на основе количества шагов и роста.
func distance(steps int, height float64) float64 {
	distanceKm := ((stepLengthCoefficient * height) * float64(steps)) / mInKm
	return distanceKm
}

// meanSpeed вычисляет среднюю скорость пользователя на основе количества шагов, роста и продолжительности тренировки.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distanceKm := distance(steps, height)
	hours := duration.Hours()
	if hours <= 0 {
		return 0
	}
	averageSpeed := distanceKm / hours
	return averageSpeed
}

// TrainingInfo формирует строку с информацией о тренировке (тип, продолжительность, дистанция, скорость, калории).
func TrainingInfo(data string, weight, height float64) (string, error) {
	step, training, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64

	switch strings.ToLower(training) {
	case "бег":
		calories, err = RunningSpentCalories(step, weight, height, duration)
	case "ходьба":
		calories, err = WalkingSpentCalories(step, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if err != nil {
		return "", errors.New("unable to calculate calories burned")
	}

	distance := distance(step, height)
	averageSpeed := meanSpeed(step, height, duration)

	return fmt.Sprintf("Тип тренировки: %s\n"+
		"Длительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\n"+
		"Сожгли калорий: %.2f\n",
		training,
		duration.Hours(),
		distance,
		averageSpeed,
		calories), nil
}

// RunningSpentCalories вычисляет количество сожженных калорий во время бега.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("неверное количество шагов")
	}
	if weight <= 0 {
		return 0, errors.New("некорректный вес")
	}
	if height <= 0 {
		return 0, errors.New("некорректный рост")
	}
	if duration <= 0 {
		return 0, errors.New("некорректное время тренировки")
	}
	averageSpeed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	calories := (weight * averageSpeed * minutes) / minInH
	return calories, nil
}

// WalkingSpentCalories вычисляет количество сожженных калорий во время ходьбы.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("неверное количество шагов")
	}
	if weight <= 0 {
		return 0, errors.New("некорректный вес")
	}
	if height <= 0 {
		return 0, errors.New("некорректный рост")
	}
	if duration <= 0 {
		return 0, errors.New("некорректное время тренировки")
	}
	minutes := duration.Minutes()
	speed := meanSpeed(steps, height, duration)
	calories := ((weight * speed * minutes) / minInH) * walkingCaloriesCoefficient

	return calories, nil
}
