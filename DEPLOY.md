# Ubuntu serverga deploy

`mms` aliasi `ubuntu` foydalanuvchisi bilan ulanadi. Server `x86_64`, shuning uchun `DEPLOY_ARCH := amd64`. Nginx backendga `127.0.0.1:8080` orqali ulanadi. Go faqat lokal build uchun kerak; serverda tayyor Linux binari ishlaydi.

## 1. Server papkasini tayyorlash

Serverga kiring:

```sh
ssh mms
```

Serverda:

```sh
sudo install -d -o ubuntu -g ubuntu -m 755 /var/www/ticket-backend
mkdir -p /var/www/ticket-backend/storage/logs
mkdir -p /var/www/ticket-backend/storage/framework/cache
mkdir -p /var/www/ticket-backend/storage/framework/sessions
mkdir -p /var/www/ticket-backend/storage/app/public
exit
```

## 2. Birinchi buildni yuklash

Kompyuteringizda backend papkasidan:

```sh
cd /Users/qodirovfaxriddin/Projects/go/ticket-backend
make build-linux
make upload
```

`make upload` binar, `public/`, `resources/` va `.env.example`ni yuklaydi. `.env` va `storage/` serverda saqlanadi. Binar avval vaqtinchalik nom bilan yuklanadi, to‘liq ko‘chirilgach almashtiriladi. Xizmatni qayta ishga tushirmaguncha ishlab turgan jarayon eski binardan foydalanadi. Bir vaqtda bitta deploy bajaring.

## 3. Serverdagi .env

```sh
ssh mms
cd /var/www/ticket-backend
cp -n .env.example .env
chmod 600 .env
nano .env
```

Quyidagi qiymatlarni sozlang; DB ma’lumotlari serverdagi PostgreSQL bazangizniki bo‘lsin:

```dotenv
APP_ENV=production
APP_DEBUG=false
APP_URL=https://bilet.amusoft.uz
APP_HOST=127.0.0.1
APP_PORT=8080
LOG_LEVEL=info

DB_CONNECTION=postgres
DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=YOUR_DATABASE
DB_USERNAME=YOUR_USERNAME
DB_PASSWORD=YOUR_PASSWORD

TELEGRAM_BOT_TOKEN=YOUR_BOT_TOKEN
TELEGRAM_INIT_DATA_TTL=1h
```

Birinchi o‘rnatishda `APP_KEY` bo‘sh bo‘lsa, yarating:

```sh
./ticket-backend artisan key:generate
```

Production tasdig‘i chiqsa, tasdiqlang. Mavjud `APP_KEY`ni keyingi deploylarda o‘zgartirmang.

Migratsiyalarni ko‘rib, kerakli bazaga qo‘llang:

```sh
./ticket-backend artisan migrate:status
./ticket-backend artisan migrate
```

Deploy migratsiyalarni avtomatik bajarmaydi.

## 4. systemd xizmati

Serverda:

```sh
sudo nano /etc/systemd/system/ticket-backend.service
```

Lokal `deploy/ticket-backend.service` fayli mazmunini joylashtiring:

```ini
[Unit]
Description=Ticket Goravel Backend
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User=ubuntu
Group=ubuntu
WorkingDirectory=/var/www/ticket-backend
ExecStart=/var/www/ticket-backend/ticket-backend
Restart=on-failure
RestartSec=5
TimeoutStopSec=30

[Install]
WantedBy=multi-user.target
```

Ishga tushiring va server yoqilganda avtomatik boshlanishini yoqing:

```sh
sudo systemctl daemon-reload
sudo systemctl enable --now ticket-backend.service
sudo systemctl status ticket-backend.service --no-pager
curl --fail http://127.0.0.1:8080/health
curl --fail http://127.0.0.1:8080/ready
```

`/health` jarayonni, `/ready` esa bazaga ulanish va zarur jadvallarni tekshiradi. Jarayon chiqargan loglar uchun `sudo journalctl -u ticket-backend.service -n 100 -f`; Goravel fayl loglari `storage/logs/` ichida.

## 5. Keyingi deploylar

Kompyuteringizda backend papkasidan:

```sh
make deploy
```

Ketma-ket Linux build qilinadi, fayllar yuklanadi va `sudo systemctl restart ticket-backend.service` bajariladi. Sudo paroli so‘ralsa, server foydalanuvchisining parolini kiriting. Birinchi deploydan oldin yuqoridagi `.env` va systemd sozlamalari tayyor bo‘lsin.

Qo‘shimcha buyruqlar:

```sh
make status
make logs
make restart
```

`deploy.mk` Git’dan chiqarilgan. Boshqa kompyuterda `make setup-deploy` orqali yarating. Undagi papka yoki xizmat nomini o‘zgartirsangiz, systemd faylidagi tegishli qiymatlarni ham moslang.

Goravel binari bilan birga kerakli kataloglar ham yuklanadi: [Goravel compile](https://docs.goravel.dev/getting-started/compile.html). Xizmatning avtomatik qayta ishga tushishi: [systemd.service](https://manpages.ubuntu.com/manpages/noble/man5/systemd.service.5.html).
