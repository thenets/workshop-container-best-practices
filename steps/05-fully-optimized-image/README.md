# 05-full-optimized-image

## How to build

## How to run

```bash
docker run -it --rm $(pwd | rev | cut -d"/" -f1 | rev)
docker run -it --rm --cap-drop CHOWN $(pwd | rev | cut -d"/" -f1 | rev)
```
