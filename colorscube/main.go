package main

import (
    "fmt"
    "time"
    "github.com/go-gl/mathgl/mgl32"
    "golang.org/x/mobile/app"
    "golang.org/x/mobile/event"
    "golang.org/x/mobile/exp/app/debug"
    "golang.org/x/mobile/geom"
    "golang.org/x/mobile/gl"
    "bytes"
    "encoding/binary"
    _ "image/png"
    "io/ioutil"
    "golang.org/x/mobile/asset"
    "golang.org/x/mobile/exp/gl/glutil"
)

type Shape struct {
    buf     gl.Buffer
    texture gl.Texture
}

type Shader struct {
    program      gl.Program
    vertCoord    gl.Attrib
    projection   gl.Uniform
    view         gl.Uniform
    model        gl.Uniform
}

type Engine struct {
    shader   Shader
    shape    Shape
    touchLoc geom.Point
    started  time.Time
}

func (e *Engine) Start() {
    var err error

    e.shader.program, err = LoadProgram("shader.v.glsl", "shader.f.glsl")
    if err != nil {
        panic(fmt.Sprintln("LoadProgram failed:", err))
    }

    e.shape.buf = gl.CreateBuffer()
    gl.BindBuffer(gl.ARRAY_BUFFER, e.shape.buf)
    fmt.Println(EncodeObject(cubeData))
    gl.BufferData(gl.ARRAY_BUFFER, EncodeObject(cubeData), gl.STATIC_DRAW)

    e.shader.vertCoord = gl.GetAttribLocation(e.shader.program, "vertCoord")

    e.shader.projection = gl.GetUniformLocation(e.shader.program, "projection")
    e.shader.view = gl.GetUniformLocation(e.shader.program, "view")
    e.shader.model = gl.GetUniformLocation(e.shader.program, "model")

    e.started = time.Now()
}

func (e *Engine) Stop() {
    gl.DeleteProgram(e.shader.program)
    gl.DeleteBuffer(e.shape.buf)
}

func (e *Engine) Config(new, old event.Config) {
    e.touchLoc = geom.Point{new.Width / 2, new.Height / 2}
}

func (e *Engine) Touch(t event.Touch, c event.Config) {
    e.touchLoc = t.Loc
}

func (e *Engine) Draw(c event.Config) {

    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LESS)

    gl.ClearColor(0, 0, 0, 1)
    gl.Clear(gl.COLOR_BUFFER_BIT)
    gl.Clear(gl.DEPTH_BUFFER_BIT)

    gl.UseProgram(e.shader.program)

    m := mgl32.Perspective(0.785, float32(c.Width/c.Height), 0.1, 10.0)
    //fmt.Println(m,"------")
    gl.UniformMatrix4fv(e.shader.projection, m[:])

    eye := mgl32.Vec3{3, 3, 3}
    center := mgl32.Vec3{0, 0, 0}
    up := mgl32.Vec3{0, 1, 0}

    m = mgl32.LookAtV(eye, center, up)
    //fmt.Println(m,"+++++++")
    gl.UniformMatrix4fv(e.shader.view, m[:])

    m = mgl32.HomogRotate3D(float32(e.touchLoc.X*6.28/c.Width), mgl32.Vec3{0, 1, 0})
    fmt.Println(m,"=======")
    gl.UniformMatrix4fv(e.shader.model, m[:])

    gl.BindBuffer(gl.ARRAY_BUFFER, e.shape.buf)

    coordsPerVertex := 3
    texCoordsPerVertex := 2
    vertexCount := len(cubeData) / (coordsPerVertex + texCoordsPerVertex)

    gl.EnableVertexAttribArray(e.shader.vertCoord)
    gl.VertexAttribPointer(e.shader.vertCoord, coordsPerVertex, gl.FLOAT, false, 20, 0) // 4 bytes in float, 5 values per vertex

    gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)

    gl.DisableVertexAttribArray(e.shader.vertCoord)

    debug.DrawFPS(c)
}

var cubeData = []float32{
    //  X, Y, Z, U, V
    // Bottom
    -1.0, -1.0, -1.0, 0.0, 0.0,
    1.0, -1.0, -1.0, 1.0, 0.0,
    -1.0, -1.0, 1.0, 0.0, 1.0,
    1.0, -1.0, -1.0, 1.0, 0.0,
    1.0, -1.0, 1.0, 1.0, 1.0,
    -1.0, -1.0, 1.0, 0.0, 1.0,

    // Top
    -1.0, 1.0, -1.0, 0.0, 0.0,
    -1.0, 1.0, 1.0, 0.0, 1.0,
    1.0, 1.0, -1.0, 1.0, 0.0,
    1.0, 1.0, -1.0, 1.0, 0.0,
    -1.0, 1.0, 1.0, 0.0, 1.0,
    1.0, 1.0, 1.0, 1.0, 1.0,

    // Front
    -1.0, -1.0, 1.0, 1.0, 0.0,
    1.0, -1.0, 1.0, 0.0, 0.0,
    -1.0, 1.0, 1.0, 1.0, 1.0,
    1.0, -1.0, 1.0, 0.0, 0.0,
    1.0, 1.0, 1.0, 0.0, 1.0,
    -1.0, 1.0, 1.0, 1.0, 1.0,

    // Back
    -1.0, -1.0, -1.0, 0.0, 0.0,
    -1.0, 1.0, -1.0, 0.0, 1.0,
    1.0, -1.0, -1.0, 1.0, 0.0,
    1.0, -1.0, -1.0, 1.0, 0.0,
    -1.0, 1.0, -1.0, 0.0, 1.0,
    1.0, 1.0, -1.0, 1.0, 1.0,

    // Left
    -1.0, -1.0, 1.0, 0.0, 1.0,
    -1.0, 1.0, -1.0, 1.0, 0.0,
    -1.0, -1.0, -1.0, 0.0, 0.0,
    -1.0, -1.0, 1.0, 0.0, 1.0,
    -1.0, 1.0, 1.0, 1.0, 1.0,
    -1.0, 1.0, -1.0, 1.0, 0.0,

    // Right
    1.0, -1.0, 1.0, 1.0, 1.0,
    1.0, -1.0, -1.0, 1.0, 0.0,
    1.0, 1.0, -1.0, 0.0, 0.0,
    1.0, -1.0, 1.0, 1.0, 1.0,
    1.0, 1.0, -1.0, 0.0, 0.0,
    1.0, 1.0, 1.0, 0.0, 1.0,
}

func main() {
    e := Engine{}
    app.Run(app.Callbacks{
        Start:  e.Start,
        Stop:   e.Stop,
        Draw:   e.Draw,
        Touch:  e.Touch,
        Config: e.Config,
    })
}



// EncodeObject converts float32 vertices into a LittleEndian byte array.
func EncodeObject(vertices ...[]float32) []byte {
    buf := bytes.Buffer{}
    for _, v := range vertices {
        err := binary.Write(&buf, binary.LittleEndian, v)
        if err != nil {
            panic(fmt.Sprintln("binary.Write failed:", err))
        }
    }

    return buf.Bytes()
}

func loadAsset(name string) ([]byte, error) {
    f, err := asset.Open(name)
    if err != nil {
        return nil, err
    }
    return ioutil.ReadAll(f)
}

// LoadProgram reads shader sources from the asset repository, compiles, and
// links them into a program.
func LoadProgram(vertexAsset, fragmentAsset string) (p gl.Program, err error) {
    vertexSrc, err := loadAsset(vertexAsset)
    if err != nil {
        return
    }

    fragmentSrc, err := loadAsset(fragmentAsset)
    if err != nil {
        return
    }

    p, err = glutil.CreateProgram(string(vertexSrc), string(fragmentSrc))
    return
}