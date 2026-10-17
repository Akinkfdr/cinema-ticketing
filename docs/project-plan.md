# Sinema bileti rezervasyon projesi

## Amaç ve kapsam

Engineering Design dersi için kullanıcıların film, salon ve seans seçip güncel
koltuk durumunu görebildiği bir uygulama geliştiriyoruz. Hedef akış: süreli koltuk
ayırma, mock ödeme, her bilet için QR üretimi ve görevlinin bileti yalnız bir kez
kullanılabilecek şekilde doğrulaması. Gerçek ödeme bu kapsamda yoktur.

Bu belge kullanıcının belirttiği teknik kararları esas alır. Ekli proje raporu
bağlamdır; rapordaki Node.js/Socket.IO önerileri seçilmiş teknoloji değildir.
Rapordaki 10 dakika, 6 koltuk ve giriş zaman aralığı gibi politikalar henüz
uygulanmış veya kesinleştirilmiş kabul edilmez.

## Mimari ve veri sorumlulukları

- Go: HTTP API, iş kuralları ve ileride WebSocket yayınları. Ekip Go ve Python
  biliyor; backend öğrenme hedefi doğrultusunda Go seçildi. Başlangıçta standart
  `net/http` yeterli; web framework, veri erişimi ve WebSocket kütüphaneleri
  ilgili aşamada değerlendirilecek. Go seçimi çift satışı kendiliğinden önlemez.
- PostgreSQL: kalıcı katalog, siparişler ve biletler. Biletlerde
  `UNIQUE(screening_id, seat_id)` kalıcı çift satışı engelleyecek.
- Redis: atomik ve süreli geçici rezervasyon; sahiplik tokenı ve TTL.
- WebSocket: yalnız sunucunun onayladığı koltuk değişikliklerini iletecek.
  Tek başına çift satışı önlemez. Bağlantı sonrası güncel durum HTTP üzerinden
  yeniden alınabilmeli. Socket.IO ile native WebSocket aynı protokol değildir.
- Docker Compose: tek backend, PostgreSQL ve Redis ile yerel geliştirme.
- React + TypeScript: sonraki aşamada değerlendirilecek frontend.

Koltuk durumu fiziksel koltuğa değil **seans + koltuk** çiftine aittir.
Bir koltuk farklı seanslarda ayrı ayrı satılabilir. Geçici durumlar AVAILABLE,
HELD ve SOLD olarak düşünülebilir; kullanılmış bir biletin koltuğu satılmış kalır.

Redis ve PostgreSQL ortak transaction paylaşmaz. Satın alma aşamasından önce
TTL sınırı ve eşzamanlı işlemler için koordinasyon tasarlanacak; rapordaki seans
satırını kilitleme yaklaşımı değerlendirilecek. Kalıcı benzersizlik kısıtı her
durumda korunacak. QR doğrulama da atomik bir durum değişimiyle yalnız bir kez
başarılı olmalı. Fiyat ve yetki kararları sunucuda verilecek.

## Geliştirme sırası ve kontrol noktaları

1. **Ortak başlangıç ortamı (mevcut aşama):** modüler Go, `GET /health`, Compose,
   örnek ortam ayarları ve Windows yönergeleri. Çıktı: API 200 ve
   `{"status":"ok"}` döndürür; üç ekip üyesi README adımlarıyla servisleri
   çalıştırır. Herkes endpoint yanıtını ve servis sağlıklarını doğrular.
2. **Veri temeli:** PostgreSQL bağlantısı, migration, film/salon/koltuk/seans
   modelleri ve küçük seed. Çıktı: aynı koltuk farklı seanslarda sorgulanabilir.
3. **Katalog ve kimlik:** listeleme, seans koltuk görünümü, müşteri/görevli
   yetkileri. Çıktı: kullanıcı seçtiği seansın koltuklarını görür.
4. **Tek koltuk rezervasyonu:** Redis atomik kilit, TTL, token kontrollü bırakma
   ve veri depoları arasında koordinasyon. Çıktı: eşzamanlı iki istekte tek başarı;
   eski token yeni rezervasyonu silemez.
5. **Mock satın alma:** transaction, benzersizlik kısıtı ve idempotency.
   Çıktı: eşzamanlı satın almada en fazla bir bilet; tekrar istekte tek sipariş.
6. **QR bilet ve doğrulama:** imzalı token ve görevli yetkisi.
   Çıktı: iki eşzamanlı doğrulamadan yalnız biri başarılı olur.
7. **Canlı durum ve frontend:** WebSocket, yeniden bağlanınca snapshot,
   React/TypeScript ekranları. Çıktı: iki istemci tutarlı durumu gösterir.
8. **Genişletme ve demo:** çoklu seçim gerekirse atomik eklenir; yarış, TTL,
   bağlantı kaybı ve yetki testleri yapılır; temiz kurulum ve ders demosu hazırlanır.

Her aşamada ilgili davranış kontrol edilerek ilerlenir. Sabit takvim yerine
çalışan küçük çıktılar kullanılır.

README sırası ortam kurulumudur; yukarıdaki sıra özellik geliştirme planıdır.
Yeni aşamaya geçmeden amacı ve araç seçimini birlikte değerlendiririz. PDF ve
önceki taslaklar zorunlu kısıt değildir; kararlar ekip ihtiyaçlarına göre değişebilir.

## Bu aşamanın sınırı

`/health` API sürecinin çalıştığını gösterir; veritabanı veya Redis bağlantısını
test etmez. Compose bu iki servis için ayrı healthcheck çalıştırır ve backend'i
onlar hazır olunca başlatır. Backend'e verilen bağlantı ortam değişkenleri ileride
kullanılacaktır; henüz veri erişimi yoktur. Redis verisi geçicidir; PostgreSQL
named volume ile korunur. Şema, migration, seed, rezervasyon, ödeme, QR,
WebSocket ve frontend bu görevde uygulanmaz.
