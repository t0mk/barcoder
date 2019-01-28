# Barcoder

.. is a CLI tool written in Go which creates PDf file with Finnish payment barcodes. It transforms YAML 

```yaml
---
- name: Normaalihintainen kausi
  iban: FI3557700520275493 
  amount: 45
  ref: 4220161

- name: Normaalikausi + jäsenmaksu (jos maksat molemmat kerralla)
  iban: FI3557700520275493 
  amount: 55
  ref: 4320168
[...]
```
into pdf so it looks like

![pdf](/pdf.png)

See 
- example input template in [salsadeleste.yml](templates/salsadeleste.yml).
- example output PDF in [example.pdf](example.pdf).


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


