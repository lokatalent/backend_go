package repository

// Repositories holds all available database tables for easy reuse.
type Repositories struct {
	User           UserRepository
	Storage        StorageRepository
	Commission     CommissionRepository
	ServicePricing ServicePricingRepository
	Booking        BookingRepository
	Notification   NotificationRepository
	Payment        PaymentRepository
}
