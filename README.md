Etapy wykonania workflow

- Wykorzystanie akcji `docker/setup-qemu-action` w celu emulacji architektury ARM64 na maszynie wirtualnej GitHub Actions (która domyślnie działa na architekturze x86_64).
- Użycie `docker/setup-buildx-action`, który aktywuje silnik Docker Buildx, wspierający budowanie manifestów wieloarchitekturalnych.
- Proces budowania został podzielony na dwa etapy. Najpierw budowany jest lokalny obraz testowy (tylko dla jednej architektury w celu oszczędności czasu). Następnie Trivy skanuje ten obraz. Flaga `exit-code: '1'` powoduje, że jeśli skaner wykryje podatności o statusie `HIGH` lub `CRITICAL`, krok ten zwraca błąd i przerywa cały potok. Dzięki temu docelowy krok `Build and push multi-arch` w ogóle się nie wykona, gwarantując bezpieczne i wolne od krytycznych wad obrazy w rejestrze `ghcr.io`.
- Strategia tagiwania:
    1. **`latest` / `edge`**: Przypisywany dla najnowszych zmian w gałęzi `main`

    2. **`sha-XXXXXXX`**: Tag oparty na skróconym commicie Git. Możemy dokładnie zidentyfikować commit, z którego powstał dany kontener.

    3. **`v*.*.*` (SemVer)**: W przypadku opublikowania tagu Git (np. `v1.0.0`), obraz otrzymuje wersję zgodną ze specyfikacją _Semantic Versioning_.
    - Uzasadnienie: strategia tagowania jest zgodna z powszechnymi praktykami rekomendowanymi przez specyfikację Open Container Initiative

        steps: - name: Checkout repo
        uses: actions/checkout@v6 # Konfiguracja QEMU dla wsparcia wieloarchitekturalności (ARM64 i AMD64) - name: Set up QEMU
        uses: docker/setup-qemu-action@v4

                - name: Set up Docker Buildx
                  uses: docker/setup-buildx-action@v4

                - name: Login to DockerHub
                  uses: docker/login-action@v4
                  with:
                      username: ${{ vars.DOCKERHUB_USERNAME }}
                      password: ${{ secrets.DOCKERHUB_TOKEN }}

                - name: Log in to GHCR
                  uses: docker/login-action@v3
                  with:
                      registry: ${{ env.REGISTRY_GHCR }}
                      username: ${{ github.actor }}
                      password: ${{ secrets.GHCR_TOKEN }}

                - name: Extract Docker metadata
                  id: meta
                  uses: docker/metadata-action@v6
                  with:
                      images: ${{ env.REGISTRY_GHCR }}/${{ github.repository_owner }}/${{ env.IMAGE_NAME }}
                      tags: |
                          type=semver,pattern={{version}}
                          type=sha,format=short
                          type=ref,event=branch
                          type=edge,branch=master

                - name: Build local image for testing
                  uses: docker/build-push-action@v6
                  with:
                      context: .
                      load: true # Ładuje obraz do lokalnego demona Dockera zamiast wysyłać do rejestru
                      tags: myapp:test
                      cache-from: type=registry,ref=${{ env.REGISTRY_DOCKERHUB }}/${{ vars.DOCKERHUB_USERNAME }}/${{ env.CACHE_IMAGE_NAME }}:cache

                - name: Run Trivy vulnerability scanner
                  uses: aquasecurity/trivy-action@master
                  with:
                      image-ref: "myapp:test"
                      format: "table"
                      exit-code: "1" # Przerywa potok, jeśli zostaną znalezione zdefiniowane podatności
                      ignore-unfixed: true
                      severity: "CRITICAL,HIGH"

                - name: Build and push multi-arch image
                  uses: docker/build-push-action@v6
                  with:
                      context: .
                      platforms: linux/amd64,linux/arm64
                      push: true
                      tags: ${{ steps.meta.outputs.tags }}
                      labels: ${{ steps.meta.outputs.labels }}
                      # Konfiguracja cache typu registry w trybie max
                      cache-from: type=registry,ref=${{ env.REGISTRY_DOCKERHUB }}/${{ vars.DOCKERHUB_USERNAME }}/${{ env.CACHE_IMAGE_NAME }}:cache
                      cache-to: type=registry,ref=${{ env.REGISTRY_DOCKERHUB }}/${{ vars.DOCKERHUB_USERNAME }}/${{ env.CACHE_IMAGE_NAME }}:cache,mode=max
