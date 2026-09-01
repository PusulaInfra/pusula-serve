# PusulaInfra Projesine Katkı Sağlama Rehberi

PusulaInfra'ya katkıda bulunmayı düşündüğünüz için teşekkür ederiz! Açık kaynak ekosistemini güçlendirmek ve kurumsal LLM altyapı standartlarını belirlemek için birlikte çalışıyoruz.

## Geliştirme Süreci
1. Bu repoyu **Fork** edin.
2. Kendi branch'inizi oluşturun (`git checkout -b feature/yeni-ozellik`).
3. Kodlarınızı yazın ve test edin (`go test ./...`).
4. Değişikliklerinizi commit edin (`git commit -m 'feat: yeni özellik eklendi'`).
5. Branch'inizi push edin (`git push origin feature/yeni-ozellik`).
6. Bir **Pull Request (PR)** açın.

## Kod Standartları
- Kodunuzun Go 1.22+ sürümleriyle uyumlu olduğundan emin olun.
- `make all` komutunu çalıştırarak derleme ve testlerin geçtiğini doğrulayın.
