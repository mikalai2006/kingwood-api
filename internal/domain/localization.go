package domain

var (
	PushTitle             = "Уведомление"
	NewMessageTitle       = "Новое сообщение"
	NewMessage            = "Добавлено сообщение для заказа: №%d - %s (объект %s)"
	NewOrderTitle         = "Новое изделие"
	NewOrder              = "Добавлено изделие: №%d - %s (объект %s)"
	CreateTaskWorkerTitle = "Новое задание"
	CreateTaskWorker      = "Добавлено задание: %s, заказ №%d - %s (объект %s)"
	PatchTaskWorkerTitle  = "Изменение задания"
	PatchTaskWorker       = "Изменено задание: %s, заказ №%d - %s (объект %s)"
	DeleteTaskWorkerTitle = "Удаление задания"
	DeleteTaskWorker      = "Удалено задание: %s, заказ №%d - %s (объект %s)"
	CreatePayTitle        = "Движение по счету за %s"
	CreatePay             = "%s пополнил ваш счет платежом - %s в размере %d ₽ за период %s"
	CreatePayAdmin        = "%s пополнил счет %s платежом - %s в размере %d ₽ за период %s"
	PatchPayTitle         = "Изменение счета за %s"
	PatchPay              = "%s изменил ваш счет %s(%d ₽) на %s(%d ₽) за период %s"
	PatchPayAdmin         = "%s изменил счет %s: %s(%d ₽) на %s(%d ₽) за период %s"
	PatchWorkTimeTitle    = "Изменение рабочей сессии"
	PatchWorkTime         = "%s изменил вашу рабочую сессию за %s: старые данные: с %s по %s(%s), новые данные: с %s по %s(%s)"
	PatchWorkTimeAdmin    = "%s изменил рабочую сессию для %s за %s: старые данные: с %s по %s(%s), новые данные: с %s по %s(%s)"
)
