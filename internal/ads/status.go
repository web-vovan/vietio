package ads

const STATUS_ACTIVE = 1
const STATUS_USER_DELETED = 2
const STATUS_EXPIRED = 3
const STATUS_SOLD = 4

func getTextStatus(codeStatus int) string {
	switch codeStatus {
	case STATUS_ACTIVE:
		return "active"
	case STATUS_USER_DELETED:
		return "deleted"
	case STATUS_EXPIRED:
		return "expired"
	case STATUS_SOLD:
		return "sold"
	default:
		return ""
	}
}
