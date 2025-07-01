# Sea Catering API

Sea Catering Backend API is a backend application for a catering service built with **Go** using **Clean Architecture**. This application provides services for managing meal plans, subscriptions, user authentication, testimonials, and payment integration.

---

## Tech Stack

- **Language**: Go 1.23.2
- **Web Framework**: Fiber v2
- **Database**: PostgreSQL
- **Cache**: Redis
- **ORM**: GORM
- **Authentication**: JWT
- **Payment Gateway**: Midtrans
- **File Storage**: Supabase Storage
- **Containerization**: Docker & Docker Compose
- **Documentation**: Swagger

---

## Features

- 🔐 **Authentication & Authorization**: JWT-based auth with Google OAuth
- 👤 **User Management**: Registration, login, profile management
- 🍽️ **Meal Plan Management**: CRUD operations for meal plans
- 📝 **Subscription System**: Subscription management with payment integration
- 💬 **Testimonials**: User reviews and ratings
- 💳 **Payment Integration**: Midtrans payment gateway
- 📧 **Email Service**: Email notifications
- 📁 **File Upload**: Image upload to Supabase storage
- 📊 **Analytics**: Dashboard stats for admin
- 🔄 **Caching**: Redis for optimal performance

---

## Folder Structure

```
.
├── .dockerignore              # Docker ignore file
├── .env                       # Environment variables (local)
├── .env.example               # Template environment variables
├── .gitignore                 # Git ignore file
├── docker-compose.yaml        # Docker compose configuration
├── Dockerfile                 # Docker build file
├── go.mod                     # Go module file
├── go.sum                     # Go module checksums
├── Makefile                   # Build automation
├── README.md                  # Project documentation
├── cmd/                       # Application entry points
│   └── api/
│       └── main.go            # Main application entry point
├── config/                    # Configuration
│   └── config.go              # Application configuration setup
├── docs/                      # API Documentation
│   ├── docs.go                # Swagger documentation generator
│   ├── swagger.json           # Swagger JSON specification
│   └── swagger.yaml           # Swagger YAML specification
└── internal/                  # Internal application code
    ├── app/                   # Application domains
    │   ├── auth/              # Authentication domain
    │   │   ├── interface/
    │   │   │   └── rest/      # REST API handlers
    │   │   ├── repository/    # Data access layer
    │   │   └── usecase/       # Business logic layer
    │   ├── meal_plan/         # Meal Plan domain
    │   │   ├── interface/
    │   │   │   └── rest/      # REST API handlers
    │   │   ├── repository/    # Data access layer
    │   │   └── usecase/       # Business logic layer
    │   ├── subscription/      # Subscription domain
    │   │   ├── interface/
    │   │   │   └── rest/      # REST API handlers
    │   │   ├── repository/    # Data access layer
    │   │   └── usecase/       # Business logic layer
    │   ├── testimonial/       # Testimonial domain
    │   │   ├── interface/
    │   │   │   └── rest/      # REST API handlers
    │   │   ├── repository/    # Data access layer
    │   │   └── usecase/       # Business logic layer
    │   └── user/              # User domain
    │       ├── interface/
    │       │   └── rest/      # REST API handlers
    │       ├── repository/    # Data access layer
    │       └── usecase/       # Business logic layer
    ├── bootstrap/             # Application initialization
    │   └── bootstrap.go       # Dependency injection & app setup
    ├── domain/                # Domain models
    │   ├── dto/               # Data Transfer Objects
    │   │   ├── auth.go        # Authentication DTOs
    │   │   ├── meal_plan.go   # Meal Plan DTOs
    │   │   ├── subscription.go# Subscription DTOs
    │   │   ├── testimonial.go # Testimonial DTOs
    │   │   └── user.go        # User DTOs
    │   └── entity/            # Database entities
    │       ├── meal_plan.go   # Meal Plan entity
    │       ├── subscription.go# Subscription entity
    │       ├── testimonial.go # Testimonial entity
    │       └── user.go        # User entity
    ├── infra/                 # Infrastructure layer
    │   ├── email/             # Email service implementation
    │   ├── fiber/             # Fiber web framework setup
    │   ├── jwt/               # JWT token implementation
    │   ├── midtrans/          # Midtrans payment gateway
    │   ├── oauth/             # OAuth implementation (Google)
    │   ├── postgresql/        # PostgreSQL database setup
    │   │   ├── migration.go   # Database migrations
    │   │   ├── postgresql.go  # Database connection
    │   │   └── seed.go        # Database seeding
    │   ├── redis/             # Redis cache implementation
    │   ├── response/          # HTTP response utilities
    │   └── supabase/          # Supabase storage integration
    ├── middleware/            # HTTP middlewares
    │   ├── authentication.go  # Authentication middleware
    │   ├── authorization.go   # Authorization middleware
    │   └── middleware.go      # Middleware interface
    └── pkg/                   # Shared packages
        ├── helper/            # Helper utilities
        ├── limiter/           # Rate limiting
        └── scheduler/         # Background job scheduler
```

### Structure Explanation

#### `/cmd`

Application entry point. Contains the `main.go` file which initializes and runs the application.

#### `/config`

Application configuration that reads environment variables and sets up global configuration.

#### `/docs`

API documentation generated by Swagger for endpoint documentation.

#### `/internal/app`

Domain-driven design where each domain has:

- **interface/rest**: HTTP handlers for the REST API
- **repository**: Data access layer for database operations
- **usecase**: Business logic layer

#### `/internal/bootstrap`

Dependency injection and application initialization. Manages all dependencies and routing.

#### `/internal/domain`

- **dto**: Data Transfer Objects for requests/responses
- **entity**: Database models with GORM tags

#### `/internal/infra`

Infrastructure layer containing implementations for:

- **email**: Service for sending emails
- **fiber**: Web framework setup
- **jwt**: JWT token management
- **midtrans**: Payment gateway integration
- **oauth**: Google OAuth implementation
- **postgresql**: Database connection and operations
- **redis**: Caching layer
- **response**: Standardized HTTP responses
- **supabase**: File storage integration

#### `/internal/middleware`

HTTP middlewares for authentication, authorization, and more.

#### `/internal/pkg`

Shared utilities and packages that can be used across the application.

---

## Setup & Installation

### Prerequisites

- Go 1.23.2+
- Docker & Docker Compose
- PostgreSQL (if not using Docker)
- Redis (if not using Docker)

### Environment Variables

Copy the `.env.example` file to `.env` and fill it with the appropriate values:

```bash
cp .env.example .env
```

Fill in the required environment variables

### Running the Application

#### Using Docker (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd sea-catering-be

# Run with Docker Compose
make run
# or
docker-compose up

# Rebuild containers
make build
# or
docker-compose up --build
```

#### Running Locally

```bash
# Install dependencies
go mod download

# Run the application
go run cmd/api/main.go
```

### Make Commands

```bash
make run                   # Run the application with Docker
make build                 # Rebuild containers
make stop                  # Stop containers
make logs                  # View logs
make restart               # Restart containers
make down-remove-volumes   # Stop and remove volumes
```

---

## API Documentation

Once the application is running, the Swagger API documentation can be accessed at:

```
http://localhost:8080/swagger/
```

---

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/auth/google` - Google OAuth login
- `GET /api/v1/auth/google/callback` - Google OAuth callback
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Logout user

### Meal Plans

- `GET /api/v1/meal-plans/` - Get all meal plans
- `GET /api/v1/meal-plans/:id` - Get meal plan by ID
- `POST /api/v1/meal-plans/` - Create a new meal plan (Admin)

### Subscriptions

- `POST /api/v1/subscriptions/` - Create a new subscription
- `GET /api/v1/subscriptions/` - Get user subscriptions
- `PUT /api/v1/subscriptions/:id/pause` - Pause a subscription
- `DELETE /api/v1/subscriptions/:id` - Cancel a subscription

### Admin Endpoints

- `GET /api/v1/subscriptions/admin/stats/new` - New subscriptions stats
- `GET /api/v1/subscriptions/admin/stats/mrr` - Monthly recurring revenue
- `GET /api/v1/subscriptions/admin/stats/active-total` - Total active subscriptions
- `GET /api/v1/subscriptions/admin/stats/reactivations` - Reactivation stats

---

## Database

The application uses **PostgreSQL** with **GORM** as the ORM. The database will be automatically migrated and seeded when the application is first run.

### Seeded Data

The application will automatically add initial data:

- **Users**: Admin and sample users
- **Meal Plans**: Diet Plan, Protein Plan, Royal Plan, Vegan Plan
- **Testimonials**: Sample testimonials

---

## License

This project is licensed under the MIT License.
