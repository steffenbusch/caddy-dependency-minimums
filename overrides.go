package dependencyminimums

// Temporary dependency overrides that should win in the final xcaddy module
// graph until Caddy itself updates to a patched minimum.
import (
	// GO-2026-5774, GO-2026-5775, GO-2026-5777
	// Found in github.com/go-chi/chi/v5@v5.2.5. Fixed in v5.3.0.
	// Remove once Caddy or the normal plugin set already pulls in v5.3.0 or newer.
	_ "github.com/go-chi/chi/v5/middleware"

	// GO-2026-5841
	// Found in github.com/klauspost/compress@v1.18.6. Fixed in v1.18.7.
	// Remove once Caddy or the normal plugin set already pulls in v1.18.7 or newer.
	_ "github.com/klauspost/compress/s2"

	// GO-2026-5158
	// Found in go.opentelemetry.io/otel@v1.43.0. Fixed in v1.44.0.
	// Remove once Caddy or the normal plugin set already pulls in v1.44.0 or newer.
	_ "go.opentelemetry.io/otel"

	// GO-2026-6061
	// Found in google.golang.org/grpc@v1.81.0. Fixed in v1.82.1.
	// Remove once Caddy or the normal plugin set already pulls in v1.82.1 or newer.
	_ "google.golang.org/grpc"
)
