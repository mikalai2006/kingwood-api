package domain

var (
	PushTitle       = "Уведомление"
	NewMessageTitle = "Новое сообщение"
	NewMessage      = "Добавлено сообщение для заказа: №%d - %s (объект %s)"

	NewOrderTitle = "Новое изделие"
	NewOrder      = "%s добавил изделие: №%d - %s (объект %s)"

	CreateTaskWorkerTitle = "Новое задание"
	CreateTaskWorker      = "%s добавил новое задание: %s, заказ №%d - %s (объект %s)"
	CreateTaskWorkerAdmin = "%s добавил новое задание для %s: %s, заказ №%d - %s (объект %s)"

	PatchTaskWorkerTitle = "Изменение задания"
	PatchTaskWorker      = "%s изменил задание: %s, заказ №%d - %s (объект %s)"
	PatchTaskWorkerAdmin = "%s изменил задание для %s: %s, заказ №%d - %s (объект %s)"

	DeleteTaskWorkerTitle = "Удаление задания"
	DeleteTaskWorker      = "%s удалил задание: %s, заказ №%d - %s (объект %s)"
	DeleteTaskWorkerAdmin = "%s удалил задание для %s: %s, заказ №%d - %s (объект %s)"

	CreatePayTitle = "Движение по счету за %s"
	CreatePay      = "%s пополнил ваш счет платежом - %s в размере %d ₽ за период %s"
	CreatePayAdmin = "%s пополнил счет %s платежом - %s в размере %d ₽ за период %s"

	PatchPayTitle = "Изменение счета за %s"
	PatchPay      = "%s изменил ваш счет %s(%d ₽) на %s(%d ₽) за период %s"
	PatchPayAdmin = "%s изменил счет %s: %s(%d ₽) на %s(%d ₽) за период %s"

	PatchWorkHistoryTitle = "Изменение рабочей сессии"
	PatchWorkHistory      = "%s изменил вашу рабочую сессию за %s: старые данные: с %s по %s(%s)(%s), новые данные: с %s по %s(%s)(%s)"
	PatchWorkHistoryAdmin = "%s изменил рабочую сессию для %s за %s: старые данные: с %s по %s(%s)(%s), новые данные: с %s по %s(%s)(%s)"

	DeleteWorkHistoryTitle = "Удаление рабочей сессии"
	DeleteWorkHistory      = "%s удалил вашу рабочую сессию за %s для заказа №%d - %s (объект %s)"
	DeleteWorkHistoryAdmin = "%s удалил рабочую сессию для %s за %s для заказа №%d - %s (объект %s)"
)
