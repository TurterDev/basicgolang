
สร้าง Project
- go mod init github.com/TurterDev/projectname
- git init
-สร้าง file main.go
สร้าง foder project structure
	|-config
		|-config.go
	|-modules
		|-users
			|-handlers (layer 1 นอกสุด)
			|-repositories (layer 3)
			|-usecase (layer 2)
			|-users.go (later ในสุด entities)
	|-tests
	|-pkg (สำหรับเก็บ source code ที่เรียกใช้งานจากภายนอก เช่น Database)
		|-database
			|-db.go
			|-migrations (ใช้สำหรับ initial database)
		|-Logger (เก็บ logger user ทีใช้งาน)
		|-authen
		|-utils (สำหรับเก็บ function พื้นฐานที่มีการเรียกใช้งานแบบ Global)
- ติดตั้ง Go migration
	- $ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	- $ migrate -version
	- https://github.com/golang-migrate/migrate?tab=readme-ov-file
	- $ cd .\pkg\database\migrations 
	- $ migrate create -ext sql -seq <database_name>
	- คำสั่ง migrate 
		$ migrate -source file://path/to/migrations -database postgres://localhost:5432/database up

- Dotenv: https://github.com/joho/godotenv
- sqlx: https://github.com/jmoiron/sqlx
- pgx: https://github.com/jackc/pgx
- GOFiber: https://docs.gofiber.io/

- ติดตั้ง live reload
	- https://github.com/cosmtrek/air
	- install => $ curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
- สร้างไฟล์แยก env => .env.dev, .env.test, .env.prod
- แก้ไขไฟล์แยก .air.dev.toml และ .air.prod.toml

- สร้าง func envPath ที่ config.go
- สร้าง func LoadConfig ที่ config.go
-install godotenv => $ go get github.com/joho/godotenv
- แก้ไขที่ file .air.dev.toml => bin = "tmp\\main.exe .env.dev" 



//Middleware คือการตัวขั้นกลางก่อนที่ user จะเข้าไปถึง API ถ้า user เข้าไปโดยไม่ผ่าน middleware ก็จะถูก reject ออก

Step ทำ Middleware
- สร้าง folder ที่ module ชื่อว่า
middlewares
	|-middlewares.go
	|-middlewaresRepositories
		|-middlewareRepositories.go
	|-middlewareUsecases
		|-middlewareUsecases.go
	|-middlewareHandlers
		|-middlewareHandlers.go
- วิธีสร้างสร้างจากส่วนลึกสุดไปนอกสุด (Repo -> Usecase -> Handler -> Module)
-เริ่มจาก middlewareRepositories
	-สร้าง interface 
	-สร้าง struct
	-สร้าง fucn
- middlewareUsecases
	- เหมือนกันกับ repo
- middlewareHandlers
	- เหมือนกับ usecase
- ลงมือทำ COSR เริ่มที่ middlewareHandlers
	- เริ่มทำ CORS
	-สร้าง func cors
	- เริ่มใช้งานโดยเพิ่ม func middleware ที่ server/module.go 
	- เพิ่ม middleware เข้าไปที่ Initmodule
	- แก้ไขที่ server.go โดยเพิ่ม middleware
-ลงมือทำ Router Check
	- สร้าง fucn RouterCheck
	- สร้าง enum
	- นำไปใช้ที่ Server.go
- ลงมือทำ logger เพิ่มเติมโดยใช้ fiber
	-ประกาศ logger เพิ่ม
	-สร้าง fucn Logger
	-นำไปใช้งานที่ Server.go

-ส้ราง Users Module
	-