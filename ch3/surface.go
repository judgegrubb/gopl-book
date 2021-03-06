// Surface computes an SVG rendering of a 3-D surface function
package main

import (
    "fmt"
    "io"
    "log"
    "math"
    "net/http"
    "os"
)

const (
    width, height   = 600, 320              // canvas size in pixels
    cells           = 100                   // number of grid cells
    xyrange         = 30.0                  // axis ranges (-xyrange..+xyrange)
    xyscale         = width / 2 / xyrange   // pixels per x or y unit
    zscale          = height * 0.4          // pixels per z unit
    angle           = math.Pi / 6           // angle of x, y axes (=30 degrees)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30degrees), cos(30degrees)

func main() {
    if len(os.Args) > 1 && os.Args[1] == "web" {
        handler := func (w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "image/svg+xml")
            surface(w)
        }
        http.HandleFunc("/", handler)
        log.Fatal(http.ListenAndServe("localhost:8000", nil))
        return
    }
    
    surface(os.Stdout)
}

func surface(out io.Writer) {
    fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' " +
        "style='stroke: grey; fill: white; stroke-width: 0.7' " +
        "width='%d' height='%d'>", width, height)
    for i := 0; i < cells; i++ {
        for j := 0; j < cells; j++ {
            ax, ay, _ := corner(i+1, j)
            bx, by, _ := corner(i, j)
            cx, cy, _ := corner(i, j+1)
            dx, dy, c := corner(i+1, j+1)
            fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g' style='fill: #%s;'/>\n", 
                ax, ay, bx, by, cx, cy, dx, dy, c)
        }
    }
    fmt.Fprintln(out, "</svg>")
}

func corner(i, j int) (float64, float64, string) {
    // Find point (x,y) at corner of cell (i,j)
    x := xyrange * (float64(i)/cells - 0.5)
    y := xyrange * (float64(j)/cells - 0.5)
    
    // Compute surface height z.
    z := f(x, y)
    
    // based on the height calculate an applicable color
    // between red and blue
    c := fmt.Sprintf("%x", int(math.Abs(z) * 255))
    if len(c) == 1 {
        c = "0" + c
    }
    if z < 0 {
        c = "0000" + c
    } else {
        c += "0000"
    }
    
    // Project (x,y,z) isometrically onto 2-D SVG canvas (sx, sy).
    sx := width / 2 + (x - y) * cos30 * xyscale
    sy := height / 2 + (x + y) * sin30 * xyscale - z * zscale
    return sx, sy, c
}

func f(x, y float64) float64  {
    r := math.Hypot(x, y) // distance from (0,0)
    return math.Sin(r) / r
}