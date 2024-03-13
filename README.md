
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
	- module
		|-users
			|-usersHandlers
				|-usersHandler.go
			|-usersRepositiries
				|-usersRepository.go
			|-usersUsecases
				|-usersUsecase.go
		|-users.go
- ประกาศ inteface, strut, constructor (usersRepository.go)
- ประกาศ inteface, strut, constructor (usersUsecase.go)
- ประกาศ inteface, strut, constructor (usersHandler.go)

-ทำ fucn sign-up (สร้าง service สำหรับ sign-up) ที่ไฟล์ users.go
	- สร้าง strut User
	- สร้าง strut UserRegister
	- สร้าง fucn สำหรับ hashpassword (BcryptHashing)
	- สร้าง func สำหรับตรวจสอบ format email

- สร้าง pettern แบบ design flegtory
	-สร้าง foder (usersPetterns) 
		|-users
			|- usersPetterns
				|-insertUser.go
	-สร้าง interface (IInertUser) โดย return Interface
	- เพิ้ม UserPassport ที่ไฟล์ users.go\
	- เพิ่ม UserToken strut
	- สร้าง private struct
	- เขียน Query SQL สำหรับ Insert
- นำ pettern ไปใช้นำ repo
	- IUsersRepository เพิ่ม

- เพิ่ม function ที่ interface IUserUsecase
	- Hashing password
	- insert User
- เพิ่ม function ที่ interface IUsersHandler
	- SignUpCustomer
		- สร่้างenum
		-Requst Body parser
		-email validation
		-Insert
- เอา register ไปใช้งาน ไปที่ module.go
	- ประกาศ func UsersModule() จะทำคล้ายๆกับ MonitorModule()
	- ทำ route ให้เป็น /v1/users/signup 
	- import UsersModule() ไปที่ไฟล์ server.go

ทำระบบ Sign-In
	- สร้าง struct ที่ไฟล์ users.go (UserCredentail) และ (UserCredentailCheck)
	- usersRepository.go สร้าง function FindOneUserByEmail 
- ไล่ทำจากข้างในไปข้างนอก Repo -> Usecase -> Hnadler -> Module
- ที่ UsersUsecase
	- สร้าง func GetPassport
	- Find user ที่ UsersUsecase
	- Set passport ที่ UsersUsecase
	- Conmpare password ที่ UsersUsecase
- ที่ usersHandler
	- สร้าง func SignIn
	- import module
- ติดตั้ง JWT
	-https://github.com/golang-jwt/jwt
	- go get -u github.com/golang-jwt/jwt/v5
- basicgolangauth
	|- สร้างไฟล์ auth.go สำหรับ generate token
- สร้าง struct UserClaims ที่ไฟล์ users.go
- NewBasicgolangAuth
- สร้าง func คำนวณวันหมดอายุ token (jwtTimeDurationCal)
- สร้าง func สำหรับ refresh token (jwtTimeRepeatAdapter)
- สร้าง func newAccessToken , newRefreshToken
- สร้าง func SignToken

- Generate Refresh Token
- เอา pagekage ที่เขียนมาใช้ ที่ usersUsecase
	- Sign Token

- Oauth
	- สร้าง func InsertOauth ที่ usersRepo
	- เรียก func ไปใช้งานที่ usersUsecase

- Refresh Token
	- สร้าง struct UserFreshCredential
	- สร้าง func สำหรับ checktoken ที่ auth.go (ParseToken)
	- check error
	- check type payload
- Repeat Token (เพื่อเปลี่ยน refresh token ใหม่)
	- สร้าง func RepeatToken
- เรียกใช้งาน Refresh Token ที่ usersRepo (จะคล้ายกับที่เขียน Signin นิดหน่อย โดยเปลี่ยนจาก insert เป็น update)
	- สร้าง func FindOneOauth ที่ usersRepo
	- สร้าง struct Oauth ที่ users.go
	- สร้าง func UpdateOauth
- Check User สำหรับ Refresh Token
	- Parse token
	- Check Oauth
	- Find Profile
	- สร้าง func UpdateOauth
	- เพิ่ม func Refresh ที่ interface
- จัดการที่ส่วน usersHandlers.go
	- จะคล้ายกับ SignIn (copy ได้เลย)
	- สร้าง errCode เพิ่ม
	- เพิ่ม func ที่ interface
	- เพิ่ม Path ในส่วนของ module
- ทำ SignOut ที่ usersRepo
	- สร้าง func DeleteOauth
	- สร้าง func DeleteOauth ต่อที่ usersUsecase
	- สร้าง func ที่ usesHandler
	- ไปที่ entities (users.go) สร้าง struct UserRemoveCredential เพิ่ม
	- เพิ่ม Path ในส่วนของ module
- ทำเพิ่ม admin
	- สร้าง func Admin ที่ InsertUser.go
	- ใช้ query เหมือน customer
	- สร้าง func InsertAdmin ที่ usersUsecase
	- สร้าง func SignUpAdmin ที่ usersHandler
	- เพิ่ม middleware ที่ module ที่ไฟล์ module
- ทำ Admin Token
	- ที่ไฟล์ Auth.go สร้าง strut โดยครอบ basicAuth
	- สร้าง interface เพิ่ม
	- สร้าง func ในการ signtoken
	- เพิ่ม func ในการ ParseToken
	- สร้าง Service สำหรับ Generate Admin key ที่ usersHandler
	- สร้าง func newAdminToken
	- เพิ่ม Path GenerateAdminToken ที่ module
- ทำ middleware Authentication
	- สร้าง func FindAccessToken ที่ middlewareRepo
	- สร้าง func FindAccessToken ที่ middlewareUsecase
	- เพิม JwtAuth() ที่ interface middlewareHandler
	- เพิม ErrCode JwtAuthErr
	- สร้าง func JwtAuth()
	- ที่ module เพิ่ม router.Get("/secret",m.mid.JwtAuth(), handler.GenerateAdminToken)
- ทำ GetUserProfile
	- สร้้าง func GetUserProfile ที่ userUsecase
	- สร้าง func GetUserProfile ที่ userHandler
- ทำ Params Check Middleware
	- สร้าง func ParamsCheck ที่ middlewareHandler
- ทำ Role-Based Authorization
	- สร้าง func FindRole ที่ middlewareRepo ([]*middlewares.Role คือ slide ของ pointer role)
	- ประกาศ type Role ที่ middlewares.go
	- สร้าง func FindRole ที่ middlewareUsecase
	- สร้าง func Authorize ที่ middlewareHandler (fucn Authorize จะทำเมื่อ JWTAuth เสร็จเรียบร้อยแล้ว) ตรง .(int) คือการแปลง type และ ...int คือการรับ type เดียวกันแบบไม่จำกัด
	สร้าง func แปลงเลขฐาน 2 ที่ utils สร้างไฟล์ converter.go
	- เรียกใช้งานที่ module.go




