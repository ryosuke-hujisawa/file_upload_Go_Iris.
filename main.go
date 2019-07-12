// main.go

package main

import (
    "os"
    "io"
    "github.com/kataras/iris"
)

const uploadsDir = "./public/uploads/"

func main() {
    app := iris.New()

    // Register templates
    app.RegisterView(iris.HTML("./views", ".html"))

    // Make the /public route path to statically serve the ./public/... contents
    app.StaticWeb("/public", "./public")

    // Render the actual form
    // GET: http://localhost:8080
    app.Get("/", func(ctx iris.Context) {
        ctx.View("upload.html")
    })

    // Upload the file to the server
    // POST: http://localhost:8080/upload
    app.Post("/upload", iris.LimitRequestBodySize(10<<20), func(ctx iris.Context) {
        // Get the file from the dropzone request
        file, info, err := ctx.FormFile("file")
        if err != nil {
            ctx.StatusCode(iris.StatusInternalServerError)
            ctx.Application().Logger().Warnf("Error while uploading: %v", err.Error())
            return
        }

        defer file.Close()
        fname := info.Filename

        // Create a file with the same name
        // assuming that you have a folder named 'uploads'
        out, err := os.OpenFile(uploadsDir+fname,
            os.O_WRONLY|os.O_CREATE, 0666)

        if err != nil {
            ctx.StatusCode(iris.StatusInternalServerError)
            ctx.Application().Logger().Warnf("Error while preparing the new file: %v", err.Error())
            return
        }
        defer out.Close()

        io.Copy(out, file)
    })

    // Start the server at http://localhost:8080
    app.Run(iris.Addr(":8080"))
}