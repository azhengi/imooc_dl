/*! https://mths.be/base64 v1.0.0 by @mathias | MIT license */
;(function(root) {

	// Detect free variables `exports`.
	var freeExports = typeof exports == 'object' && exports;

	// Detect free variable `module`.
	var freeModule = typeof module == 'object' && module &&
		module.exports == freeExports && module;

	// Detect free variable `global`, from Node.js or Browserified code, and use
	// it as `root`.
	var freeGlobal = typeof global == 'object' && global;
	if (freeGlobal.global === freeGlobal || freeGlobal.window === freeGlobal) {
		root = freeGlobal;
	}

	/*--------------------------------------------------------------------------*/

	var InvalidCharacterError = function(message) {
		this.message = message;
	};
	InvalidCharacterError.prototype = new Error;
	InvalidCharacterError.prototype.name = 'InvalidCharacterError';

	var error = function(message) {
		// Note: the error messages used throughout this file match those used by
		// the native `atob`/`btoa` implementation in Chromium.
		throw new InvalidCharacterError(message);
	};

	var TABLE = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
	// http://whatwg.org/html/common-microsyntaxes.html#space-character
	var REGEX_SPACE_CHARACTERS = /[\t\n\f\r ]/g;

	// `decode` is designed to be fully compatible with `atob` as described in the
	// HTML Standard. http://whatwg.org/html/webappapis.html#dom-windowbase64-atob
	// The optimized base64-decoding algorithm used is based on @atk’s excellent
	// implementation. https://gist.github.com/atk/1020396
	var decode = function(input) {
		input = String(input)
			.replace(REGEX_SPACE_CHARACTERS, '');
		var length = input.length;
		if (length % 4 == 0) {
			input = input.replace(/==?$/, '');
			length = input.length;
		}
		if (
			length % 4 == 1 ||
			// http://whatwg.org/C#alphanumeric-ascii-characters
			/[^+a-zA-Z0-9/]/.test(input)
		) {
			error(
				'Invalid character: the string to be decoded is not correctly encoded.'
			);
		}
		var bitCounter = 0;
		var bitStorage;
		var buffer;
		var output = '';
		var position = -1;
		while (++position < length) {
			buffer = TABLE.indexOf(input.charAt(position));
			bitStorage = bitCounter % 4 ? bitStorage * 64 + buffer : buffer;
			// Unless this is the first of a group of 4 characters…
			if (bitCounter++ % 4) {
				// …convert the first 8 bits to a single ASCII character.
				output += String.fromCharCode(
					0xFF & bitStorage >> (-2 * bitCounter & 6)
				);
			}
		}
		return output;
	};

	// `encode` is designed to be fully compatible with `btoa` as described in the
	// HTML Standard: http://whatwg.org/html/webappapis.html#dom-windowbase64-btoa
	var encode = function(input) {
		input = String(input);
		if (/[^\0-\xFF]/.test(input)) {
			// Note: no need to special-case astral symbols here, as surrogates are
			// matched, and the input is supposed to only contain ASCII anyway.
			error(
				'The string to be encoded contains characters outside of the ' +
				'Latin1 range.'
			);
		}
		var padding = input.length % 3;
		var output = '';
		var position = -1;
		var a;
		var b;
		var c;
		var buffer;
		// Make sure any padding is handled outside of the loop.
		var length = input.length - padding;

		while (++position < length) {
			// Read three bytes, i.e. 24 bits.
			a = input.charCodeAt(position) << 16;
			b = input.charCodeAt(++position) << 8;
			c = input.charCodeAt(++position);
			buffer = a + b + c;
			// Turn the 24 bits into four chunks of 6 bits each, and append the
			// matching character for each of them to the output.
			output += (
				TABLE.charAt(buffer >> 18 & 0x3F) +
				TABLE.charAt(buffer >> 12 & 0x3F) +
				TABLE.charAt(buffer >> 6 & 0x3F) +
				TABLE.charAt(buffer & 0x3F)
			);
		}

		if (padding == 2) {
			a = input.charCodeAt(position) << 8;
			b = input.charCodeAt(++position);
			buffer = a + b;
			output += (
				TABLE.charAt(buffer >> 10) +
				TABLE.charAt((buffer >> 4) & 0x3F) +
				TABLE.charAt((buffer << 2) & 0x3F) +
				'='
			);
		} else if (padding == 1) {
			buffer = input.charCodeAt(position);
			output += (
				TABLE.charAt(buffer >> 2) +
				TABLE.charAt((buffer << 4) & 0x3F) +
				'=='
			);
		}

		return output;
	};

	var base64 = {
		'encode': encode,
		'decode': decode,
		'version': '1.0.0'
	};

	// Some AMD build optimizers, like r.js, check for specific condition patterns
	// like the following:
	if (
		typeof define == 'function' &&
		typeof define.amd == 'object' &&
		define.amd
	) {
		define(function() {
			return base64;
		});
	}	else if (freeExports && !freeExports.nodeType) {
		if (freeModule) { // in Node.js or RingoJS v0.8.0+
			freeModule.exports = base64;
		} else { // in Narwhal or RingoJS v0.7.0-
			for (var key in base64) {
				base64.hasOwnProperty(key) && (freeExports[key] = base64[key]);
			}
		}
	} else { // in Rhino or a web browser
		root.base64 = base64;
	}

}(this));


var result, j
var K = function(a) {
    function e(b) {
        N = b;
        G = Array(N);
        for (b = 0; b < G.length; b++)
            G[b] = 0;
        new h;
        E = new h;
        E.digits[0] = 1
    }
    function h(b) {
        this.digits = "boolean" == typeof b && 1 == b ? null : G.slice(0);
        this.isNeg = !1
    }
    function u(b) {
        var c = new h(!0);
        c.digits = b.digits.slice(0);
        c.isNeg = b.isNeg;
        return c
    }
    function k(b) {
        for (var c = new h, r = b.length, d = 0; 0 < r; r -= 4,
        ++d) {
            for (var a = c.digits, O = d, g = b.substr(Math.max(r - 4, 0), Math.min(r, 4)), e = 0, f = Math.min(g.length, 4), l = 0; l < f; ++l) {
                e <<= 4;
                var n = g.charCodeAt(l);
                e |= 48 <= n && 57 >= n ? n - 48 : 65 <= n && 90 >= n ? 10 + n - 65 : 97 <= n && 122 >= n ? 10 + n - 97 : 0
            }
            a[O] = e
        }
        return c
    }
    function w(b, c) {
        if (b.isNeg != c.isNeg) {
            c.isNeg = !c.isNeg;
            var r = l(b, c);
            c.isNeg = !c.isNeg
        } else {
            r = new h;
            for (var d = 0, a = 0; a < b.digits.length; ++a)
                d = b.digits[a] + c.digits[a] + d,
                r.digits[a] = d & 65535,
                d = Number(65536 <= d);
            r.isNeg = b.isNeg
        }
        return r
    }
    function l(b, c) {
        if (b.isNeg != c.isNeg) {
            c.isNeg = !c.isNeg;
            var r = w(b, c);
            c.isNeg = !c.isNeg
        } else {
            r = new h;
            for (var a, m = a = 0; m < b.digits.length; ++m)
                a = b.digits[m] - c.digits[m] + a,
                r.digits[m] = a & 65535,
                0 > r.digits[m] && (r.digits[m] += 65536),
                a = 0 - Number(0 > a);
            if (-1 == a) {
                for (m = a = 0; m < b.digits.length; ++m)
                    a = 0 - r.digits[m] + a,
                    r.digits[m] = a & 65535,
                    0 > r.digits[m] && (r.digits[m] += 65536),
                    a = 0 - Number(0 > a);
                r.isNeg = !b.isNeg
            } else
                r.isNeg = b.isNeg
        }
        return r
    }
    function n(b) {
        for (var c = b.digits.length - 1; 0 < c && 0 == b.digits[c]; )
            --c;
        return c
    }
    function v(b) {
        var c = n(b);
        b = b.digits[c];
        c = 16 * (c + 1);
        var a;
        for (a = c; a > c - 16 && 0 == (b & 32768); --a)
            b <<= 1;
        return a
    }
    function t(b, c) {
        for (var a = new h, d, m = n(b), e = n(c), g, f = 0; f <= e; ++f) {
            d = 0;
            g = f;
            for (j = 0; j <= m; ++j,
            ++g)
                d = a.digits[g] + b.digits[j] * c.digits[f] + d,
                a.digits[g] = d & 65535,
                d >>>= 16;
            a.digits[f + m + 1] = d
        }
        a.isNeg = b.isNeg != c.isNeg;
        return a
    }
    function p(b, c, a, d, m) {
        for (m = Math.min(c + m, b.length); c < m; ++c,
        ++d)
            a[d] = b[c]
    }
    function y(b, c) {
        var a = Math.floor(c / 16)
          , d = new h;
        p(b.digits, 0, d.digits, a, d.digits.length - a);
        c %= 16;
        a = 16 - c;
        for (var m = d.digits.length - 1, e = m - 1; 0 < m; --m,
        --e)
            d.digits[m] = d.digits[m] << c & 65535 | (d.digits[e] & P[c]) >>> a;
        d.digits[0] = d.digits[m] << c & 65535;
        d.isNeg = b.isNeg;
        return d
    }
    function L(b, a) {
        var c = Math.floor(a / 16)
          , d = new h;
        p(b.digits, c, d.digits, 0, b.digits.length - c);
        a %= 16;
        c = 16 - a;
        for (var e = 0, f = e + 1; e < d.digits.length - 1; ++e,
        ++f)
            d.digits[e] = d.digits[e] >>> a | (d.digits[f] & Q[a]) << c;
        d.digits[d.digits.length - 1] >>>= a;
        d.isNeg = b.isNeg;
        return d
    }
    function C(a, c) {
        var b = new h;
        p(a.digits, 0, b.digits, c, b.digits.length - c);
        return b
    }
    function x(a, c) {
        var b = new h;
        p(a.digits, c, b.digits, 0, b.digits.length - c);
        return b
    }
    function D(a, c) {
        var b = new h;
        p(a.digits, 0, b.digits, 0, c);
        return b
    }
    function M(a, c) {
        if (a.isNeg != c.isNeg)
            return 1 - 2 * Number(a.isNeg);
        for (var b = a.digits.length - 1; 0 <= b; --b)
            if (a.digits[b] != c.digits[b])
                return a.isNeg ? 1 - 2 * Number(a.digits[b] > c.digits[b]) : 1 - 2 * Number(a.digits[b] < c.digits[b]);
        return 0
    }
    function F(a) {
        this.modulus = u(a);
        this.k = n(this.modulus) + 1;
        a = new h;
        a.digits[2 * this.k] = 1;
        var c = this.modulus
          , b = v(a)
          , d = v(c)
          , e = c.isNeg;
        if (b < d)
            if (a.isNeg) {
                var f = u(E);
                f.isNeg = !c.isNeg;
                a.isNeg = !1;
                c.isNeg = !1;
                var g = l(c, a);
                a.isNeg = !0;
                c.isNeg = e
            } else
                f = new h,
                g = u(a);
        else {
            f = new h;
            g = a;
            for (var q = Math.ceil(d / 16) - 1, k = 0; 32768 > c.digits[q]; )
                c = y(c, 1),
                ++k,
                ++d,
                q = Math.ceil(d / 16) - 1;
            g = y(g, k);
            b = Math.ceil((b + k) / 16) - 1;
            for (d = C(c, b - q); -1 != M(g, d); )
                ++f.digits[b - q],
                g = l(g, d);
            for (; b > q; --b) {
                d = b >= g.digits.length ? 0 : g.digits[b];
                var p = b - 1 >= g.digits.length ? 0 : g.digits[b - 1]
                  , t = b - 2 >= g.digits.length ? 0 : g.digits[b - 2]
                  , B = q >= c.digits.length ? 0 : c.digits[q]
                  , z = q - 1 >= c.digits.length ? 0 : c.digits[q - 1];
                f.digits[b - q - 1] = d == B ? 65535 : Math.floor((65536 * d + p) / B);
                for (var A = f.digits[b - q - 1] * (65536 * B + z), x = 4294967296 * d + (65536 * p + t); A > x; )
                    --f.digits[b - q - 1],
                    A = f.digits[b - q - 1] * (65536 * B | z),
                    x = 4294967296 * d + (65536 * p + t);
                t = d = C(c, b - q - 1);
                B = f.digits[b - q - 1];
                result = new h;
                p = n(t);
                for (z = A = 0; z <= p; ++z)
                    A = result.digits[z] + t.digits[z] * B + A,
                    result.digits[z] = A & 65535,
                    A >>>= 16;
                result.digits[1 + p] = A;
                g = l(g, result);
                g.isNeg && (g = w(g, d),
                --f.digits[b - q - 1])
            }
            g = L(g, k);
            f.isNeg = a.isNeg != e;
            a.isNeg && (f = e ? w(f, E) : l(f, E),
            c = L(c, k),
            g = l(c, g));
            0 == g.digits[0] && 0 == n(g) && (g.isNeg = !1)
        }
        a = [f, g];
        this.mu = a[0];
        this.bkplus1 = new h;
        this.bkplus1.digits[this.k + 1] = 1;
        this.modulo = H;
        this.multiplyMod = I;
        this.powMod = J
    }
    function H(a) {
        var b = x(a, this.k - 1);
        b = t(b, this.mu);
        b = x(b, this.k + 1);
        a = D(a, this.k + 1);
        b = t(b, this.modulus);
        b = D(b, this.k + 1);
        a = l(a, b);
        a.isNeg && (a = w(a, this.bkplus1));
        for (b = 0 <= M(a, this.modulus); b; )
            a = l(a, this.modulus),
            b = 0 <= M(a, this.modulus);
        return a
    }
    function I(a, c) {
        a = t(a, c);
        return this.modulo(a)
    }
    function J(a, c) {
        var b = new h;
        for (b.digits[0] = 1; ; ) {
            0 != (c.digits[0] & 1) && (b = this.multiplyMod(b, a));
            c = L(c, 1);
            if (0 == c.digits[0] && 0 == n(c))
                break;
            a = this.multiplyMod(a, a)
        }
        return b
    }
    function K(a) {
        this.e = k("10001");
        this.d = k("");
        this.m = k(a);
        this.chunkSize = 128;
        this.radix = 16;
        this.barrett = new F(this.m)
    }
    var N, G, E;
    e(20);
    (function(a) {
        var b = new h;
        b.isNeg = 0 > a;
        a = Math.abs(a);
        for (var f = 0; 0 < a; )
            b.digits[f++] = a & 65535,
            a >>= 16;
        return b
    }
    )(1E15);
    var P = [0, 32768, 49152, 57344, 61440, 63488, 64512, 65024, 65280, 65408, 65472, 65504, 65520, 65528, 65532, 65534, 65535]
      , Q = [0, 1, 3, 7, 15, 31, 63, 127, 255, 511, 1023, 2047, 4095, 8191, 16383, 32767, 65535];
    e(131);
    return function(a) {
        var b = [], e = a.length, d, m = "", l = new K(f);
        e > l.chunkSize - 11 && (e = l.chunkSize - 11);
        var g = 0;
        for (d = e - 1; g < e; )
            b[d] = a.charCodeAt(g),
            g++,
            d--;
        for (d = l.chunkSize - e % l.chunkSize; 0 < d; ) {
            for (a = Math.floor(256 * Math.random()); !a; )
                a = Math.floor(256 * Math.random());
            b[g] = a;
            g++;
            d--
        }
        b[e] = 0;
        b[l.chunkSize - 2] = 2;
        b[l.chunkSize - 1] = 0;
        e = b.length;
        for (g = 0; g < e; g += l.chunkSize) {
            var q = new h;
            d = 0;
            for (a = g; a < g + l.chunkSize; ++d)
                q.digits[d] = b[a++],
                q.digits[d] += b[a++] << 8;
            d = l.barrett.powMod(q, l.e);
            q = "";
            for (a = n(d); -1 < a; --a) {
                var k = d.digits[a];
                var p = String.fromCharCode(k & 255);
                k = String.fromCharCode(k >>> 8 & 255) + p;
                q += k
            }
            d = q;
            m += d
        }
        return m
    }(a)
}

function n(t, e) {
    function r(t, e) {
        var r = "";
        if ("object" == typeof t)
            for (var n = 0; n < t.length; n++)
                r += String.fromCharCode(t[n]);
        t = r || t;
        for (var i, o, a = new Uint8Array(t.length), s = e.length, n = 0; n < t.length; n++)
            o = n % s,
            i = t[n],
            i = i.toString().charCodeAt(0),
            a[n] = i ^ e.charCodeAt(o);
        return a
    }

    function n(t) {
        var e = "";
        if ("object" == typeof t)
            for (var r = 0; r < t.length; r++)
                e += String.fromCharCode(t[r]);
        t = e || t;
        var n = new Uint8Array(t.length);
        for (r = 0; r < t.length; r++)
            n[r] = t[r].toString().charCodeAt(0);
        var i, o, r = 0;
        for (r = 0; r < n.length; r++)
            0 != (i = n[r] % 3) && r + i < n.length && (o = n[r + 1],
                n[r + 1] = n[r + i],
                n[r + i] = o,
                r = r + i + 1);
        return n
    }

    function i(t) {
        var e = "";
        if ("object" == typeof t)
            for (var r = 0; r < t.length; r++)
                e += String.fromCharCode(t[r]);
        t = e || t;
        var n = new Uint8Array(t.length);
        for (r = 0; r < t.length; r++)
            n[r] = t[r].toString().charCodeAt(0);
        var r = 0,
            i = 0,
            o = 0,
            a = 0;
        for (r = 0; r < n.length; r++)
            o = n[r] % 2,
            o && r++,
            a++;
        var s = new Uint8Array(a);
        for (r = 0; r < n.length; r++)
            o = n[r] % 2,
            s[i++] = o ? n[r++] : n[r];
        return s
    }

    function o(t, e) {
        var r = 0,
            n = 0,
            i = 0,
            o = 0,
            a = "";
        if ("object" == typeof t)
            for (var r = 0; r < t.length; r++)
                a += String.fromCharCode(t[r]);
        t = a || t;
        var s = new Uint8Array(t.length);
        for (r = 0; r < t.length; r++)
            s[r] = t[r].toString().charCodeAt(0);

        for (r = 0; r < t.length; r++) {
            if (0 != (o = s[r] % 5) && 1 != o && r + o < s.length && (i = s[r + 1],
                    n = r + 2,
                    s[r + 1] = s[r + o],
                    s[o + r] = i,
                    (r = r + o + 1) - 2 > n)) {
                for (; n < r - 2; n++) {
                    s[n] = s[n] ^ e.charCodeAt(n % e.length)
                }
            }
        }
        for (r = 0; r < t.length; r++)
            s[r] = s[r] ^ e.charCodeAt(r % e.length);

        return s
    }
    for (var a = {
            data: {
                info: t
            }
        }, s = {
            q: r,
            h: n,
            m: i,
            k: o
        }, l = a.data.info, u = l.substring(l.length - 4).split(""), c = 0; c < u.length; c++)
        u[c] = u[c].toString().charCodeAt(0) % 4;
    u.reverse();
    for (var d = [], c = 0; c < u.length; c++)
        d.push(l.substring(u[c] + 1, u[c] + 2)),
        l = l.substring(0, u[c] + 1) + l.substring(u[c] + 2);
    a.data.encrypt_table = d,
        a.data.key_table = [];
    for (var c in a.data.encrypt_table)
        "q" != a.data.encrypt_table[c] && "k" != a.data.encrypt_table[c] || (a.data.key_table.push(l.substring(l.length - 12)),
            l = l.substring(0, l.length - 12));
    a.data.key_table.reverse(),
        a.data.info = l;
    var f = new Array(-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 62, -1, -1, -1, 63, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, -1, -1, -1, -1, -1, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, -1, -1, -1, -1, -1, -1, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, -1, -1, -1, -1, -1);
    a.data.info = function(t) {
        var e, r, n, i, o, a, s;
        for (a = t.length,
            o = 0,
            s = ""; o < a;) {
            do {
                e = f[255 & t.charCodeAt(o++)]
            } while (o < a && -1 == e);
            if (-1 == e)
                break;
            do {
                r = f[255 & t.charCodeAt(o++)]
            } while (o < a && -1 == r);
            if (-1 == r)
                break;
            s += String.fromCharCode(e << 2 | (48 & r) >> 4);
            do {
                if (61 == (n = 255 & t.charCodeAt(o++)))
                    return s;
                n = f[n]
            } while (o < a && -1 == n);
            if (-1 == n)
                break;
            s += String.fromCharCode((15 & r) << 4 | (60 & n) >> 2);
            do {
                if (61 == (i = 255 & t.charCodeAt(o++)))
                    return s;
                i = f[i]
            } while (o < a && -1 == i);
            if (-1 == i)
                break;
            s += String.fromCharCode((3 & n) << 6 | i)
        }
        return s
    }(a.data.info);
    for (var c in a.data.encrypt_table) {
        var h = a.data.encrypt_table[c];
        if ("q" == h || "k" == h) {
            var p = a.data.key_table.pop();
            a.data.info = s[a.data.encrypt_table[c]](a.data.info, p)
        } else{
            a.data.info = s[a.data.encrypt_table[c]](a.data.info)
            }
    }
    if (e)
        return a.data.info;
    var g = "";
    for (c = 0; c < a.data.info.length; c++)
        g += String.fromCharCode(a.data.info[c]);
    return g
};


var f = "DBCEA86ACD310CC0ED8A56D9E3C3CFE26951E8A6C0AC103419B43617C410B0537B13E7D145AB007E61BB39CB66854A4AA9BABD108BD93212376CD9A61A03B80B03D54D760F8FD317C784AE1B8489A2D3890ABCF3F73946EEBF7CF433BB4C53526DE29F4CFECF07F3C95CF2A95BF140EE605F695FF0889EECFD3F6808C85254B1";

function handleDecrypt(info, e) {
    return n(info, e);
}

function handleCryptPassword(pw) {
    var r = base64.encode(K(pw));
    
    return r;
};