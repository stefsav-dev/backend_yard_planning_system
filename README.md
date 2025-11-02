# Backend Yard Planning System 

Sebuah backend service untuk manajemen penempatan kontainer di pelabuhan dengan sistem perencanaan yard yang cerdas.

📋 Deskripsi Sistem
Sistem Yard Planning merupakan backend service yang mengelola penempatan kontainer di area pelabuhan. Sistem ini membagi yard menjadi beberapa block, dan setiap block memiliki struktur spasial (slot, row, tier) untuk menentukan posisi unik kontainer.

🚀 Fitur Utama
Yard & Block Management - Mengelola data yard dan block storage

Smart Container Placement - Sistem saran posisi otomatis berdasarkan spesifikasi kontainer

Container Tracking - Melacak penempatan dan pengambilan kontainer

Redis Caching - Optimasi performa dengan caching multi-layer

RESTful API - API lengkap untuk integrasi dengan frontend/system lain

🛠️ Teknologi yang Digunakan
Backend: Go (Golang) dengan Fiber Framework

Database: PostgreSQL dengan GORM ORM

Caching: Redis untuk optimasi performa

Validation: Go Playground Validator

Environment: Configurable environment variables

📦 Struktur Database
Tabel-tabel Utama
yards - Data yard/pelabuhan

blocks - Block storage dalam yard

yard_plans - Rencana penempatan berdasarkan spesifikasi kontainer

containers - Data kontainer dan posisinya

