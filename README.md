# MRegeln

Regelset für klassisches Pen and Paper Rollenspiel. Die Regeln liegen als Schnipsel in Markdown Dateien vor.
Das Repository stellt ein Kommandozeilentool in golang bereit das aus diesen Schnipseln mit Hilfe einer Skriptdatei
ein Regelwerk in LaTeX generiert. Somit können modular Regeln für das jeweilige Setting zusammengestellt werden.

Das Kommandozeilentool beeinhaltet weiter nützliche Funktionen zum Umgang mit dem Regelwerk

Ein Grund das Regelwerk mittels Git zu versionieren besteht darin es weiter entwickeln zu können und gegebenenfalls
"breaking changes" veröfflichen zu können (z.B eine neuen Version in einem neuen "branch").

## Usage

* Kommandozeilentool ist für den eigenen Gebrauch erstellt! (primär generieren eines LaTeX Dokuments)
* Golang Kenntnisse erforderlich um das Kommandozeiletool zu kompilieren. 
* Bei Interesse wird das Tool kompiliert bereitgestellt.
* Das kompilierte Tool stellt Hilfe mit dem `-h` oder `--help` Schalter bereit
