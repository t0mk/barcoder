# Barcoder

.. is a CLI tool written in Go which creates PDf file with Finnish payment barcodes. The PDf can be printed 


## Build 

Install Golang, adn then in the root of this repo:

```
$ go build
```


## Usage

```
./barcoder --templ templates/salsadeleste.yml --outfile "example.pdf" --date 2019-02-04
```

It will generate pdf to [example.pdf](example.pdf).

## Barcode

It generates the barcode for Finnish bank apps, which can usually be scanned with a smartphone.


