package database

import (
	"encoding/json"
	"log"
	"time"

	"app/model"

	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	var count int64

	// --- Resources ---
	db.Model(&model.Resource{}).Count(&count)
	if count == 0 {
		resources := []model.Resource{
			// Industrial Materials
			{Name: "Сталь", Description: "Конструкционная сталь высокого качества", Unit: "кг", Quantity: 1000},
			{Name: "Алюминий", Description: "Алюминиевые листы для производства", Unit: "кг", Quantity: 500},
			{Name: "Медь", Description: "Медные провода и кабели", Unit: "м", Quantity: 2000},
			{Name: "Пластик ПВХ", Description: "Поливинилхлорид для изготовления труб", Unit: "кг", Quantity: 750},
			{Name: "Стекло", Description: "Листовое стекло различной толщины", Unit: "м²", Quantity: 300},

			// Energy Resources
			{Name: "Электроэнергия", Description: "Потребление электрической энергии", Unit: "кВт·ч", Quantity: 5000},
			{Name: "Природный газ", Description: "Газ для отопления и производства", Unit: "м³", Quantity: 1200},
			{Name: "Дизельное топливо", Description: "Топливо для генераторов и техники", Unit: "л", Quantity: 800},
			{Name: "Уголь", Description: "Каменный уголь для котельной", Unit: "т", Quantity: 50},

			// Liquids and Chemicals
			{Name: "Вода", Description: "Техническая вода для производства", Unit: "л", Quantity: 10000},
			{Name: "Питьевая вода", Description: "Очищенная питьевая вода", Unit: "л", Quantity: 2000},
			{Name: "Кислота серная", Description: "Серная кислота для химических процессов", Unit: "л", Quantity: 150},
			{Name: "Щелочь натрия", Description: "Гидроксид натрия", Unit: "кг", Quantity: 200},
			{Name: "Растворитель", Description: "Органические растворители", Unit: "л", Quantity: 300},

			// Building Materials
			{Name: "Цемент", Description: "Портландцемент М400", Unit: "т", Quantity: 20},
			{Name: "Песок", Description: "Речной песок строительный", Unit: "м³", Quantity: 100},
			{Name: "Щебень", Description: "Гранитный щебень фракция 5-20мм", Unit: "м³", Quantity: 80},
			{Name: "Кирпич", Description: "Красный керамический кирпич", Unit: "шт", Quantity: 5000},
			{Name: "Арматура", Description: "Стальная арматура А500С", Unit: "т", Quantity: 15},

			// Office Supplies
			{Name: "Бумага A4", Description: "Офисная бумага белая", Unit: "пачка", Quantity: 100},
			{Name: "Картриджи", Description: "Картриджи для принтеров", Unit: "шт", Quantity: 25},
			{Name: "Канцтовары", Description: "Ручки, карандаши, скрепки", Unit: "набор", Quantity: 50},

			// Tools and Equipment
			{Name: "Сверла", Description: "Сверла по металлу различных диаметров", Unit: "шт", Quantity: 200},
			{Name: "Болты", Description: "Болты М8-М20 различной длины", Unit: "шт", Quantity: 1000},
			{Name: "Гайки", Description: "Гайки к болтам М8-М20", Unit: "шт", Quantity: 1200},
			{Name: "Шайбы", Description: "Плоские и пружинные шайбы", Unit: "шт", Quantity: 2000},

			// IT Equipment
			{Name: "Серверы", Description: "Серверное оборудование", Unit: "шт", Quantity: 5},
			{Name: "Мониторы", Description: "ЖК мониторы 24 дюйма", Unit: "шт", Quantity: 30},
			{Name: "Клавиатуры", Description: "USB клавиатуры", Unit: "шт", Quantity: 40},
			{Name: "Мыши", Description: "Оптические USB мыши", Unit: "шт", Quantity: 40},
		}

		if err := db.Create(&resources).Error; err != nil {
			log.Println("❌ Ошибка при добавлении ресурсов:", err)
		} else {
			log.Println("✅ Добавлены тестовые ресурсы")

			// Seed resource history after resources are created
			seedResourceHistory(db)
		}
	}
}

// seedResourceHistory creates sample history entries for resources
func seedResourceHistory(db *gorm.DB) {
	var count int64
	db.Model(&model.ResourceHistory{}).Count(&count)
	if count > 0 {
		return // History already exists
	}

	// Get some resources and users for history entries
	var resources []model.Resource
	var users []model.User

	db.Limit(10).Find(&resources)
	db.Find(&users)

	if len(resources) == 0 || len(users) == 0 {
		log.Println("⚠️ Не удалось найти ресурсы или пользователей для создания истории")
		return
	}

	// Create history entries
	historyEntries := []model.ResourceHistory{}

	// Initial creation entries for first few resources
	for i, resource := range resources[:5] {
		newDataJSON, _ := json.Marshal(resource)

		historyEntries = append(historyEntries, model.ResourceHistory{
			ResourceID:  resource.ID,
			Action:      "CREATE",
			UserID:      users[i%len(users)].ID,
			OldData:     "",
			NewData:     string(newDataJSON),
			Timestamp:   time.Now().Add(-time.Duration(i*24) * time.Hour),
			Description: "Ресурс '" + resource.Name + "' создан",
		})
	}

	// Some update entries
	for i, resource := range resources[1:4] {
		oldResource := resource
		oldResource.Quantity = resource.Quantity / 2

		oldDataJSON, _ := json.Marshal(oldResource)
		newDataJSON, _ := json.Marshal(resource)

		historyEntries = append(historyEntries, model.ResourceHistory{
			ResourceID:  resource.ID,
			Action:      "UPDATE",
			UserID:      users[(i+1)%len(users)].ID,
			OldData:     string(oldDataJSON),
			NewData:     string(newDataJSON),
			Timestamp:   time.Now().Add(-time.Duration(i*12) * time.Hour),
			Description: "Количество ресурса '" + resource.Name + "' обновлено",
		})
	}

	// Inventory adjustments
	for i, resource := range resources[2:6] {
		oldResource := resource
		oldResource.Quantity = resource.Quantity + 100

		oldDataJSON, _ := json.Marshal(oldResource)
		newDataJSON, _ := json.Marshal(resource)

		historyEntries = append(historyEntries, model.ResourceHistory{
			ResourceID:  resource.ID,
			Action:      "UPDATE",
			UserID:      users[(i+2)%len(users)].ID,
			OldData:     string(oldDataJSON),
			NewData:     string(newDataJSON),
			Timestamp:   time.Now().Add(-time.Duration(i*6) * time.Hour),
			Description: "Корректировка запасов ресурса '" + resource.Name + "'",
		})
	}

	if err := db.Create(&historyEntries).Error; err != nil {
		log.Println("❌ Ошибка при добавлении истории ресурсов:", err)
	} else {
		log.Printf("✅ Добавлено %d записей истории ресурсов", len(historyEntries))
	}
}
