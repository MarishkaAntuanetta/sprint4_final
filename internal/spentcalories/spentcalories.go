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
		return 0, "", 0, errors.New("неверный формат данных")
	}

	step, err1 := strconv.Atoi(strData[0])
	if err1 != nil || step <= 0 {
		return 0, "", 0, errors.New("неверное значение шага")
	}

	duration, err2 := time.ParseDuration(strData[2])
	if err2 != nil || duration <= 0 {
		return 0, "", 0, errors.New("неверное значение продолжительности")
	}

	return step, strData[1], duration, nil
}

// distance вычисляет дистанцию, пройденную пользователем на основе количества шагов и роста.
func distance(steps int, height float64) float64 {
	distanceKm := ((stepLengthCoefficient * height) * float64(steps)) / float64(mInKm)
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

	switch strings.ToLower(training) {
	case "бег":
		distance := distance(step, height)
		averageSpeed := meanSpeed(step, height, duration)
		calories, err := RunningSpentCalories(step, weight, height, duration)
		if err != nil {
			return "", errors.New("невозможно произвести подсчет калорий")
		}
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

	case "ходьба":
		distance := distance(step, height)
		averageSpeed := meanSpeed(step, height, duration)
		calories, err := WalkingSpentCalories(step, weight, height, duration)
		if err != nil {
			return "", errors.New("невозможно произвести подсчет калорий")
		}
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

	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

// RunningSpentCalories вычисляет количество сожженных калорий во время бега.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("неверные данные")
	}
	averageSpeed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	calories := (weight * averageSpeed * minutes) / minInH
	return calories, nil
}

// WalkingSpentCalories вычисляет количество сожженных калорий во время ходьбы.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("неверные данные")
	}

	minutes := duration.Minutes()
	speed := meanSpeed(steps, height, duration)
	calories := ((weight * speed * minutes) / minInH) * walkingCaloriesCoefficient

	return calories, nil
}
