import gzip
import hashlib
import math
import sys


def h(x):
    return abs(int(hashlib.sha1(str(x)).hexdigest(), 16))


def p(s):
    r = 0
    while (s & 1) != 1:
        s >>= 1
        r += 1
    return r + 1


def alpha(m):
    if m == 16:
        return 0.673
    if m == 32:
        return 0.697
    if m == 64:
        return 0.709
    return 0.7213 / (1.0 + (1.079 / m))


def hll(f, b):
    m = 1 << b
    M = [0 for _ in range(m)]
    bs = 0
    for i in range(b):
        bs |= (1 << i)

    exact = set()

    with gzip.open(f, 'r') as f:
        for line in f:
            x = h(line)
            j = x & bs
            w = x >> b
            M[j] = max(M[j], p(w))

            exact.add(int(line))

    Z = sum([1.0 / (1 << e) for e in M])
    E = (alpha(m) * m * m) / Z

    if E < (5.0 / 2.0) * m:
        V = len([e for e in M if e == 0])
        if V != 0:
            E = m * math.log(m / V, 2)

    return E, len(exact)


if __name__ == '__main__':
    f = sys.argv[1]
    b = 8

    est, ext = hll(f, b)
    print('Estimated value: %s' % est)
    print('Exact value: %s' % ext)
    print('Approximation: %.3f' % (est/ext))
