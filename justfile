# Default project binary name
binary := "cinder"

# Default recipe: Lists available options
default:
    @just --list

# Builds the interactive Go CLI utility from the src folder
build:
    @echo "🛠️ Compiling interactive Go CLI..."
    go build -o {{binary}} ./cmd/cinder
    @echo "✨ Compilation completed successfully! Executable generated: ./{{binary}}"

# Builds and runs the interactive cleanup utility
run: build
    @echo "🚀 Starting Cinder CLI..."
    ./{{binary}}

# Removes the compiled executable from the root folder
clean:
    @echo "🧹 Removing compiled executable..."
    rm -f {{binary}}
    @echo "✨ Cleanup completed!"

# Organizes and tidies Go module dependencies
tidy:
    @echo "📦 Tidying Go modules..."
    go mod tidy
    @echo "✨ Dependencies tidied!"

# Installs the compiled binary globally in ~/.local/bin
install: build
    @echo "🚀 Installing Cinder in ~/.local/bin..."
    mkdir -p ~/.local/bin
    cp {{binary}} ~/.local/bin/
    @echo "✨ Cinder successfully installed! Try running 'cinder' in your terminal."
