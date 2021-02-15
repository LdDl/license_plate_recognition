# Create folder if not exists
mkdir -p data

# license_plates weights (it was trained on Russian license plates dataset)
wget --load-cookies /tmp/cookies.txt "https://docs.google.com/uc?export=download&confirm=$(wget --quiet --save-cookies /tmp/cookies.txt --keep-session-cookies --no-check-certificate 'https://docs.google.com/uc?export=download&id=1aqlI3jDzE0QJJ5oQ_9ODlceYmZDGRKqF' -O- | sed -rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1aqlI3jDzE0QJJ5oQ_9ODlceYmZDGRKqF" -O data/license_plates_15000.weights && rm -rf /tmp/cookies.txt

# license_plates names (currently it contains only license_car, license_taxi)
wget --load-cookies /tmp/cookies.txt "https://docs.google.com/uc?export=download&confirm=$(wget --quiet --save-cookies /tmp/cookies.txt --keep-session-cookies --no-check-certificate 'https://docs.google.com/uc?export=download&id=1Keai52qz-XTDr4ruHSN8_wcN5MaMa5NM' -O- | sed -rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1Keai52qz-XTDr4ruHSN8_wcN5MaMa5NM" -O data/license_plates.names && rm -rf /tmp/cookies.txt

# license_plates inference cfg 
wget --load-cookies /tmp/cookies.txt "https://docs.google.com/uc?export=download&confirm=$(wget --quiet --save-cookies /tmp/cookies.txt --keep-session-cookies --no-check-certificate 'https://docs.google.com/uc?export=download&id=1EAmjHb99YZ2aqLE_TwJqKC_UDLEpJq8Z' -O- | sed -rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1EAmjHb99YZ2aqLE_TwJqKC_UDLEpJq8Z" -O data/license_plates_inference.cfg && rm -rf /tmp/cookies.txt
## Don't forget to add [names] to *.cfg file. It's needed for AlexeyAB's fork
sed -i -e "\$anames = ../data/license_plates.names" data/license_plates_inference.cfg

# ocr weights (YOLO v4)
wget --load-cookies /tmp/cookies.txt "https://docs.google.com/uc?export=download&confirm=$(wget --quiet --save-cookies /tmp/cookies.txt --keep-session-cookies --no-check-certificate 'https://docs.google.com/uc?export=download&id=1MZ2ii0hQmKpIcwj3Mfh5DOBnlmV7KL5T' -O- | sed -rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1MZ2ii0hQmKpIcwj3Mfh5DOBnlmV7KL5T" -O data/ocr_plates_140000.weights && rm -rf /tmp/cookies.txt

# ocr names (YOLO v4, there are 22 possible symbols in Russian license plates)
wget --load-cookies /tmp/cookies.txt "https://docs.google.com/uc?export=download&confirm=$(wget --quiet --save-cookies /tmp/cookies.txt --keep-session-cookies --no-check-certificate 'https://docs.google.com/uc?export=download&id=1654bphBaeQ6LJUZJO_NthPMNpp4oJXQA' -O- | sed -rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1654bphBaeQ6LJUZJO_NthPMNpp4oJXQA" -O data/ocr_plates.names && rm -rf /tmp/cookies.txt

# ocr inference cfg (YOLO v4)
wget --load-cookies /tmp/cookies.txt "https://docs.google.com/uc?export=download&confirm=$(wget --quiet --save-cookies /tmp/cookies.txt --keep-session-cookies --no-check-certificate 'https://docs.google.com/uc?export=download&id=1d-IdpviI8imGHJYmGz8C33KHSLBvUBZo' -O- | sed -rn 's/.*confirm=([0-9A-Za-z_]+).*/\1\n/p')&id=1d-IdpviI8imGHJYmGz8C33KHSLBvUBZo" -O data/ocr_plates_inference.cfg && rm -rf /tmp/cookies.txt
## Again: don't forget to add [names] to *.cfg file. It's needed for AlexeyAB's fork
sed -i -e "\$anames = ../data/ocr_plates.names" data/ocr_plates_inference.cfg
