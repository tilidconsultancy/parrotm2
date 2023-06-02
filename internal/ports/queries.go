package ports

import "github.com/google/uuid"

func GetByIds(ids []uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"id": map[string]interface{}{
			"$in": ids,
		},
	}
}

func GetById(id uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"id": id,
	}
}

func GetTenantByPhoneId(pid string) map[string]interface{} {
	return map[string]interface{}{
		"accountsettings.phoneid": pid,
	}
}

func GetConversationByTenantAndUser(pid string, uphone string) map[string]interface{} {
	return map[string]interface{}{
		"tenant.accountsettings.phoneid": pid,
		"user.phone":                     uphone,
	}
}
