Алгоритм:

Пока остались ребра, берем ребро из очереди, добавляем в граф, смотрим формирует ли оно 3 или 4-цикл, если да - удаляем, если нет, оставляем. И так пока все ребра не пройдем. Перед всем этим упорядочиваем ребра (сначала было рандомное перемешивание всего списка ребер, потом сортировка по длинам ребер по убыванию кусками размера w1 и рандомное перемешивание w2 раз, w1/w2 подбирались экспериментально). При этом все распараллелено.
