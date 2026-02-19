package email

import (
	"encoding/json"

	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/internal/utils"
)

type OverdueItem struct {
	Name         string
	ReturnedDate string
}

type OverdueEmailData struct {
	Total int
	More  int
	Items []OverdueItem
}

func BuildOverdueEmailData(notifications []*model.Notification, limit int) OverdueEmailData {
	total := len(notifications)

	display := notifications
	more := 0

	if total > limit {
		display = notifications[:limit]
		more = total - limit
	}

	items := make([]OverdueItem, 0, len(display))
	for _, n := range display {
		var payload model.NotificationPayload

		if err := json.Unmarshal(n.Payload, &payload); err != nil {
			continue
		}
		name, _ := payload.Extra["name"].(string)
		returnedRaw, _ := payload.Extra["returned_date"].(string)
		formattedDate := utils.FormatDate(returnedRaw)

		items = append(items, OverdueItem{
			Name:         name,
			ReturnedDate: formattedDate,
		})
	}

	return OverdueEmailData{
		Total: total,
		More:  more,
		Items: items,
	}
}
