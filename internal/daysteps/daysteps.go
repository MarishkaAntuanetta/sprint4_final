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
	// stepLength — длина одного шага в метрах.
	stepLength = 0.65
	// mInKm — количество метров в одном километре.
	mInKm = 1000
)

// parsePackage разбирает входную строку с данными о шагах и времени активности.
func parsePackage(data string) (int, time.Duration, error) {
	// Разделяем входные данные по запятой
	strData := strings.Split(data, ",")
	if len(strData) != 2 {
		return 0, 0, errors.New("invalid data format")
	}

	// Извлекаем количество шагов и продолжительность
	stepStr := strData[0]
	durationStr := strData[1]

	// Проверяем на лишние пробелы — если обрезанная версия не равна оригинальной, значит формат неверный
	if stepStr != strings.TrimSpace(stepStr) || durationStr != strings.TrimSpace(durationStr) {
		return 0, 0, errors.New("invalid data format — extra spaces")
	}

	// Очистка пробелов и знака "+" перед числом шагов
	stepStr = strings.TrimPrefix(strings.TrimSpace(stepStr), "+")
	durationStr = strings.TrimSpace(durationStr)

	// Преобразуем количество шагов в число
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		return 0, 0, errors.New("invalid data format")
	}
	if step <= 0 {
		return 0, 0, errors.New("invalid step value")
	}
	// Преобразуем продолжительность в тип time.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, errors.New("invalid data format")
	}
	if duration <= 0 {
		return 0, 0, errors.New("invalid activity duration")
	}

	return step, duration, nil
}

// DayActionInfo формирует строку с информацией о выполненной активности
func DayActionInfo(data string, weight, height float64) string {
	// Парсим данные о шагах и продолжительности
	step, duration, err := parsePackage(data)

	// Если возникла ошибка при разборе данных, логируем и возвращаем пустую строку
	if err != nil {
		log.Printf("%s\n", err.Error())
		return ""
	}

	// Рассчитываем дистанцию в метрах и километрах
	distance := float64(step) * stepLength
	distanceKm := distance / mInKm

	// Считаем потраченные калории
	caloriesBurned, err := spentcalories.WalkingSpentCalories(step, weight, height, duration)
	if err != nil {
		log.Println("unable to calculate calories burned")
		return ""
	}

	// Формируем итоговую строку с результатами
	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		step,
		distanceKm,
		caloriesBurned,
	)
}
