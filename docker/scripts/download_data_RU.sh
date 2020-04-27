#!/bin/bash

# Create folder if not exists
mkdir -p /models/data

# license_plates weights (it was trained on Russian license plates dataset)
wget -nc --load-cookies /tmp/cookies.txt "https://docs.google.com/uc?export=download&confirm=$(wget --quiet --save-cookies /tmp/cookies.txt --keep-session-cookies --no-check-certificate 'https://docs.google.com/uc?export=download&id=1aqlI3jDzE0QJJ5oQ_9ODlceYmZDGRKqF' -O- | sed -rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1aqlI3jDzE0QJJ5oQ_9ODlceYmZDGRKqF" -O /models/data/license_plates_15000.weights && rm -rf /tmp/cookies.txt

# license_plates names (currently it contains only license_car, license_taxi)
wget -nc --load-cookies /tmp/cookies.txt "https://docs.google.com/uc?export=download&confirm=$(wget --quiet --save-cookies /tmp/cookies.txt --keep-session-cookies --no-check-certificate 'https://docs.google.com/uc?export=download&id=1Keai52qz-XTDr4ruHSN8_wcN5MaMa5NM' -O- | sed -rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1Keai52qz-XTDr4ruHSN8_wcN5MaMa5NM" -O /models/data/license_plates.names && rm -rf /tmp/cookies.txt

# license_plates inference cfg 
wget -nc --load-cookies /tmp/cookies.txt "https://docs.google.com/uc?export=download&confirm=$(wget --quiet --save-cookies /tmp/cookies.txt --keep-session-cookies --no-check-certificate 'https://docs.google.com/uc?export=download&id=1EAmjHb99YZ2aqLE_TwJqKC_UDLEpJq8Z' -O- | sed -rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1EAmjHb99YZ2aqLE_TwJqKC_UDLEpJq8Z" -O /models/data/license_plates_inference.cfg && rm -rf /tmp/cookies.txt
## Don't forget to add [names] to *.cfg file. It's needed for AlexeyAB's fork
sed -i -e "\$anames = /models/data/license_plates.names" /models/data/license_plates_inference.cfg

# ocr weights
wget -nc --load-cookies /tmp/cookies.txt "https://docs.google.com/uc?export=download&confirm=$(wget --quiet --save-cookies /tmp/cookies.txt --keep-session-cookies --no-check-certificate 'https://docs.google.com/uc?export=download&id=1WQwWuYfiU8ZxFEqDoZvH2bKUBUw797U9' -O- | sed -rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1WQwWuYfiU8ZxFEqDoZvH2bKUBUw797U9" -O /models/data/ocr_plates_7000.weights && rm -rf /tmp/cookies.txt

# ocr names (there are 22 possible symbols in Russian license plates)
wget -nc --load-cookies /tmp/cookies.txt "https://docs.google.com/uc?export=download&confirm=$(wget --quiet --save-cookies /tmp/cookies.txt --keep-session-cookies --no-check-certificate 'https://docs.google.com/uc?export=download&id=1SLPw1zT_WwwN46mg2ttQv5DIl5HEZRNj' -O- | sed -rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1SLPw1zT_WwwN46mg2ttQv5DIl5HEZRNj" -O /models/data/ocr_plates.names && rm -rf /tmp/cookies.txt

# ocr inference cfg
wget -nc --load-cookies /tmp/cookies.txt "https://docs.google.com/uc?export=download&confirm=$(wget --quiet --save-cookies /tmp/cookies.txt --keep-session-cookies --no-check-certificate 'https://docs.google.com/uc?export=download&id=1KQvN6l1G0OxyJj__YXgeYyAWxlEqzvfz' -O- | sed -rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1KQvN6l1G0OxyJj__YXgeYyAWxlEqzvfz" -O /models/data/ocr_plates_inference.cfg && rm -rf /tmp/cookies.txt
## Again: don't forget to add [names] to *.cfg file. It's needed for AlexeyAB's fork
sed -i -e "\$anames = /models/data/ocr_plates.names" /models/data/ocr_plates_inference.cfg
