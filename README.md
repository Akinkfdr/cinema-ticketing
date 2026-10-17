# cinema-ticketing
Real-time cinema seat reservation and QR ticket verification

Engineering Design dersi için sinema bileti rezervasyon uygulaması.
Bu aşama yalnız başlangıç altyapısını içerir: Go, PostgreSQL ve Redis.
Proje amacı ve geliştirme sırası: [proje planı](docs/project-plan.md).
Çalışma kuralları: [AGENTS.md](AGENTS.md).

## 1. Windows hazırlığı

Docker Desktop'ı Linux containers/WSL 2 desteğiyle kurup başlatın. PowerShell'i
proje kökünde açın. Docker Desktop'ın hazır olmasını bekleyin:

```powershell
docker version
docker compose version
```

`docker version` hem Client hem Server göstermeli. Named pipe/daemon hatası
alırsanız Docker Desktop'ın ve Linux engine'in çalıştığını kontrol edin.

## 2. Yerel ayarlar

```powershell
Copy-Item .env.example .env
```

Bu komutu yalnız `.env` henüz yoksa çalıştırın; mevcut ayarlarınızı koruyun.
`.env` Git tarafından yok sayılır. Örnek parola yalnız yerel geliştirme içindir.
Portlar doluysa `.env` içindeki `BACKEND_PORT`, `POSTGRES_PORT`, `REDIS_PORT`
değerlerini değiştirin. Servisler yalnız `127.0.0.1` üzerinde yayınlanır.

## 3. Yapılandırmayı kontrol edin ve başlatın

```powershell
docker compose config --quiet
docker compose up --build -d --wait
docker compose ps
```

İlk build internetten Go derleyici ve servis imajlarını indirir, ardından Go
kodunu derler. Üç servis de `healthy` görünmelidir. PostgreSQL/Redis hazır
olduktan sonra backend başlar. Backend yalnız Go standart kütüphanesini kullanır;
henüz üçüncü taraf Go bağımlılığı yoktur. Docker ile çalıştırırken bilgisayara
ayrıca Go kurmak gerekmez.

## 4. Çalıştığını doğrulayın

```powershell
Invoke-RestMethod http://localhost:8000/health
docker compose exec postgres sh -c 'pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"'
docker compose exec redis redis-cli ping
docker compose logs --tail 50 backend postgres redis
```

Beklenen sonuçlar sırasıyla `status: ok`, `accepting connections`, `PONG`.
HTTP endpoint'i 200 ve JSON olarak `{"status":"ok"}` yanıtını döndürür. Bu aşamada otomatik Swagger arayüzü veya `/docs` adresi yoktur.
Portu değiştirdiyseniz URL'deki 8000'i de değiştirin.

`/health` yalnız API'nin çalıştığını bildirir; PostgreSQL/Redis bağlantısını
sınamaz. Backend henüz bu servislere sorgu göndermez. Compose içinden bağlantı
adresleri `postgres:5432` ve `redis:6379`, Windows üzerinden `localhost` ve
`.env` içindeki yayın portlarıdır. Container içinde `localhost` aynı container'dır.

## 5. Değiştirme ve durdurma

Kod değişikliğinden sonra `docker compose up --build -d --wait` çalıştırın.
Bu başlangıç yapılandırmasında otomatik kod yenileme yoktur.

```powershell
docker compose down
```

PostgreSQL verisi `postgres_data` named volume'ünde kalır. `down -v` bu veriyi
siler; normal durdurmada kullanmayın. Mevcut PostgreSQL volume'ünde kullanıcı,
parola veya veritabanı değişikliği `.env` değişince otomatik uygulanmaz.
Redis yeniden başladığında geçici veriler kaybolur; kalıcı bilet deposu değildir.

## Docker olmadan yalnız API'yi kontrol etme

Go 1.26.5 kuruluysa proje kökünde aşağıdaki komutları çalıştırın. Docker build
aynı Go sürümünü kullanır. Go modülünün minimum sürümü 1.26.0'dır.

```powershell
go version
Set-Location backend
go test ./...
go vet ./...
go run ./cmd/api
```

İkinci PowerShell penceresinde:

```powershell
Invoke-RestMethod http://localhost:8000/health
```

`status: ok` beklenir. Sunucuyu `Ctrl+C` ile durdurun. Backend klasöründen
proje köküne dönmek için `Set-Location ..` kullanın. Compose komutlarını kökte
çalıştırın. Yerel sunucu varsayılan olarak `127.0.0.1:8000` dinler. Port doluysa,
`go run` öncesinde `$env:HTTP_ADDR = "127.0.0.1:8001"` ayarlayabilirsiniz;
kontrol URL'sini de değiştirin. Compose `.env` dosyasını okur; `go run` bu dosyayı
otomatik okumaz. Docker içindeki sunucu `HTTP_ADDR=:8000` ile dinler.

Bu yöntem PostgreSQL ve Redis'i başlatmaz. Python, pip ve `.venv` artık bu
backend için gerekli değildir. Önceki denemeden kalan yerel `.venv` klasörü
kullanılmaz ve Git'e gönderilmez.

## Üç kişilik ekip için ilk kontrol noktası

Her ekip üyesi paylaşılan aynı proje sürümünde yukarıdaki Docker kurulumunu
uygular, üç servisin `healthy` olduğunu ve `/health` yanıtını kontrol eder.
Herkes çalıştırabildiğinde ortak ortam adımı tamamlanır. Başlangıç dosyalarını
GitHub deposunu klonlayarak alın; mevcut kopyanız varsa yerel değişikliklerinizi
koruyarak güncel sürüme geçin.

README ortamı çalıştırma sırasını, `docs/project-plan.md` ise özellik geliştirme
sırasını anlatır. Sonraki özelliğe geçmeden kapsamı ve araçları birlikte seçeriz.

## Dosya yapısı

```text
backend/
  cmd/api/main.go             # Sunucuyu başlatır
  internal/httpapi/
    router.go                # URL ve handler eşleştirmeleri
    health.go                # GET /health yanıtı
    router_test.go           # HTTP sözleşmesi ve yönlendirme testleri
  go.mod                     # Modül adı ve minimum Go sürümü
  Dockerfile                 # Derleme ve çalıştırma imajları
  .dockerignore              # Build'e gönderilecek dosyalar
docs/project-plan.md          # Amaç, kararlar, geliştirme sırası
docs/verification.md          # Çalıştırılan kontrollerin sonuçları
AGENTS.md                    # Kod üzerinde çalışırken uyulacak kurallar
compose.yaml                 # Üç servisin ortak çalışma ayarları
.env.example                 # Paylaşılabilir yerel ayar örneği
```

Koltuk kilitleme, mock ödeme, QR, WebSocket ve frontend sonraki görevlerdir.
