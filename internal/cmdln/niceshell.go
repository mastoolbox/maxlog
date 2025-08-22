package cmdln

import (
	"strings"
	"log"
)

const (
    Reset = "\033[0m"

    // Foreground colors
    Black        = "\033[30m"
    Red          = "\033[31m"
    Green        = "\033[32m"
    Yellow       = "\033[33m"
    Blue         = "\033[34m"
    Magenta      = "\033[35m"
    Cyan         = "\033[36m"
    LightGray    = "\033[37m"
    DarkGray     = "\033[90m"
    LightRed     = "\033[91m"
    LightGreen   = "\033[92m"
    LightYellow  = "\033[93m"
    LightBlue    = "\033[94m"
    LightMagenta = "\033[95m"
    LightCyan    = "\033[96m"
    White        = "\033[97m"

    // Background colors
    BGBlack      = "\033[40m"
    BGRed        = "\033[41m"
    BGGreen      = "\033[42m"
    BGYellow     = "\033[43m"
    BGBlue       = "\033[44m"
    BGMagenta    = "\033[45m"
    BGCyan       = "\033[46m"
    BGLightGray  = "\033[47m"
    BGLightBlue  = "\033[104m"
)

const SymError = 0
const SymWarn  = 1

var aSymbols = map[int]string {
    0: Red + "\uebfb" + Reset,
    1: Yellow + "\uf071" + Reset,
}

// GetSymbol retrieves a predefined symbol based on its identifier.
//
// Parameter:
//   id - An integer representing the identifier of the symbol.
//
// Returns:
//   string - The symbol associated with the given identifier.
func GetSymbol(id int) string {
    return aSymbols[id]
}

// Fatal logs a fatal error message and terminates the program.
//
// Parameters:
//   msg - The error message to log.
//   err - An optional error object to include in the log.
//
// Behavior:
//   - Constructs a formatted error message using a predefined symbol and the provided message.
//   - If the `err` parameter is not nil, includes the error object in the log.
//   - Terminates the program execution using `log.Fatal`.
func Fatal(msg string, err error) {
    text := GetSymbol(SymError) + " " + msg
    if err != nil {
        log.Fatal(text, err)
    } else {
        log.Fatal(text)
    }
}

// SetLabels processes the input text by applying color-coded labels to specific substrings.
//
// Parameters:
//   text - The input string to process.
//   tag  - An optional tag to highlight in the text.
//
// Returns:
//   string - The processed text with color-coded labels applied.
//
// Behavior:
//   - Iterates through a predefined list of replacements, where each replacement specifies:
//     - oldKey: The substring to be replaced.
//     - newKey: The replacement text to insert.
//     - labelFunc: The function used to apply the color-coded label.
//   - Applies the corresponding label function to replace substrings in the text.
//   - If the text contains "CID-CRON" or "BMXAA6372I", applies subdued formatting using the Downplay function.
func SetLabels(text, tag string) string {

    replacements := []struct {
        oldKey string
        newKey string
        labelFunc func(string, string, string) string
    }{
        {"[INFO]", "INFO", SetBlueLabel},
        {"[INFO ]", "INFO", SetBlueLabel},
        {"[AUDIT   ]", "AUDIT", SetBlueLabel},
        {"[WARN]", "WARN", SetYellowLabel},
        {"[WARN ]", "WARN", SetYellowLabel},
        {"[WARNING ]", "WARN", SetYellowLabel},
        {"[ERROR]", "ERROR", SetRedLabel},
        {"[ERROR   ]", "ERROR", SetRedLabel},
        {"[err]", "ERROR", SetRedLabel},
        {"[MXServer]", "MX", SetMagentaLabel},
        {"[MAXIMO_UI]", "UI", SetMagentaLabel},
        {"[maximo]", "MAX", SetCyanLabel},
        {"[DEBUG]", "DEBUG", SetCyanLabel},
        {"[maximo.script." + tag +"]", "Script", SetLightBlueLabel},
        {"Maximo is ready for client connections.", "Maximo is ready for client connections.", SetGreenLabel},
    }

    for _, r := range replacements {
        text = r.labelFunc(text, r.oldKey, r.newKey)
    }

    if strings.Contains(text, "CID-CRON") || strings.Contains(text, " BMXAA6372I") {
        text = Downplay(text)
    }

    if tag != "" && strings.Contains(text, tag)  && text[0] != '\t' {
        text = Highlight(text)
        text = SetGreenLabel(text, tag, tag)
    }

    return text
}

// SetLabel replaces a specific substring in the input text with a formatted label.
//
// Parameters:
//   text     - The input string to process.
//   oldKey   - The substring to be replaced.
//   newKey   - The replacement text to insert.
//   color    - The foreground color to apply to the label.
//   bgColor  - The background color to apply to the label.
//
// Returns:
//   string - The processed text with the formatted label.
func SetLabel(text, oldKey, newKey, color, bgColor string) string {
    return strings.Replace(text, oldKey, color+"\ue0b6"+Reset+bgColor+newKey+Reset+color+"\ue0b4"+Reset, 1)
}

// SetYellowLabel applies a yellow color-coded label to a specific substring in the input text.
//
// Parameters:
//   text    - The input string to process.
//   oldKey  - The substring to be replaced.
//   newKey  - The replacement text to insert.
//
// Returns:
//   string - The processed text with the yellow label.
func SetYellowLabel(text, oldKey, newKey string) string {
    return SetLabel(text, oldKey, newKey, Yellow, BGYellow)
}

// SetBlueLabel applies a blue color-coded label to a specific substring in the input text.
//
// Parameters:
//   text    - The input string to process.
//   oldKey  - The substring to be replaced.
//   newKey  - The replacement text to insert.
//
// Returns:
//   string - The processed text with the blue label.
func SetBlueLabel(text, oldKey, newKey string) string {
    return SetLabel(text, oldKey, newKey, Blue, BGBlue)
}

// SetLightBlueLabel applies a light blue color-coded label to a specific substring in the input text.
//
// Parameters:
//   text    - The input string to process.
//   oldKey  - The substring to be replaced.
//   newKey  - The replacement text to insert.
//
// Returns:
//   string - The processed text with the blue label.
func SetLightBlueLabel(text, oldKey, newKey string) string {
    return SetLabel(text, oldKey, newKey, LightBlue, BGLightBlue)
}

// SetRedLabel applies a red color-coded label to a specific substring in the input text.
//
// Parameters:
//   text    - The input string to process.
//   oldKey  - The substring to be replaced.
//   newKey  - The replacement text to insert.
//
// Returns:
//   string - The processed text with the red label.
func SetRedLabel(text, oldKey, newKey string) string {
    return SetLabel(text, oldKey, newKey, Red, BGRed)
}

// SetMagentaLabel applies a magenta color-coded label to a specific substring in the input text.
//
// Parameters:
//   text    - The input string to process.
//   oldKey  - The substring to be replaced.
//   newKey  - The replacement text to insert.
//
// Returns:
//   string - The processed text with the magenta label.
func SetMagentaLabel(text, oldKey, newKey string) string {
    return SetLabel(text, oldKey, newKey, Magenta, BGMagenta)
}

// SetCyanLabel applies a cyan color-coded label to a specific substring in the input text.
//
// Parameters:
//   text    - The input string to process.
//   oldKey  - The substring to be replaced.
//   newKey  - The replacement text to insert.
//
// Returns:
//   string - The processed text with the cyan label.
func SetCyanLabel(text, oldKey, newKey string) string {
    return SetLabel(text, oldKey, newKey, Cyan, BGCyan)
}

// SetGreenLabel applies a green color-coded label to a specific substring in the input text.
//
// Parameters:
//   text    - The input string to process.
//   oldKey  - The substring to be replaced.
//   newKey  - The replacement text to insert.
//
// Returns:
//   string - The processed text with the green label.
func SetGreenLabel(text, oldKey, newKey string) string {
    return SetLabel(text, oldKey, newKey, Green, BGGreen)
}

// Downplay applies a subdued formatting to specific substrings in the input text.
//
// Parameter:
//   text - The input string to process.
//
// Returns:
//   string - The processed text with subdued formatting.
//
// Behavior:
//   - Replaces the first occurrence of " [" in the text with a dark gray color-coded version.
//   - Appends a reset sequence to ensure proper formatting.
func Downplay(text string) (string) {
    if strings.Contains(text, " [") {
        return strings.Replace(text, " [", DarkGray + " [" , 1) + Reset
     }
     return text
}

// Highlight applies a white color-coded formatting to specific substrings in the input text.
//
// Parameter:
//   text - The input string to process.
//
// Returns:
//   string - The processed text with white color-coded formatting.
//
// Behavior:
//   - Replaces the first occurrence of " [" in the text with a white color-coded version.
//   - Appends a reset sequence to ensure proper formatting.
func Highlight(text string) (string) {
    if strings.Contains(text, " [") {
        return strings.Replace(text, " [", White + " [" , 1) + Reset
    }
    return text
}


