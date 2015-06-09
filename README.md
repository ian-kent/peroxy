Peroxy
======

Take control of your browser. MITM style.

### Huh?

A MITM proxy and script injector.

For controlling TV browsers. From the command line!

**Stating the obvious:** not for use in the wild

### What?

From the command line:

```bash
~> go get github.com/ian-kent/peroxy
~> peroxy
```

From a browser: http://your-ip-address:3123

From the command line:
```bash
~> curl http://your-ip-address:3123/!-switch?proxy=https://www.google.co.uk&url=/
~> curl http://your-ip-address:3123/!-switch?proxy=https://www.microsoft.com&url=/
~> curl http://your-ip-address:3123/!-switch?proxy=https://www.mozilla.org&url=/
```

Evaluate some javascript:
```bash
~> curl http://your-ip-address:3123/!-switch?eval=alert\(\'hi\'\)\;
```

### Licence

Copyright ©‎ 2015, Ian Kent (http://iankent.uk).

Released under MIT license, see [LICENSE](LICENSE.md) for details.
