# @rixycf dnsmasq container for adblock
FROM alpine:3.6

LABEL maintainer="rixycf kasnake1013@gmail.com"

RUN apk add --no-cache dnsmasq curl bind-tools && \
    curl -sL "https://warui.intaa.net/adhosts/hosts.txt" | \
    awk 'NR > 1 {print "address=/"$2"/"$1}' >> /etc/dnsmasq.adblock.conf && \
    echo "domain-needed" >> /etc/dnsmasq.conf && \
    echo "bogus-priv" >> /etc/dnsmasq.conf && \
    echo "no-resolv" >> /etc/dnsmasq.conf && \
    echo "no-poll" >> /etc/dnsmasq.conf && \
    echo "server=8.8.8.8" >> /etc/dnsmasq.conf && \
    echo "server=8.8.4.4" >> /etc/dnsmasq.conf && \
    echo "conf-file=/etc/dnsmasq.adblock.conf" >> /etc/dnsmasq.conf

HEALTHCHECK --interval=30s --timeout=30s --start-period=10s --retries=2 \
    CMD dig @localhost 0.r.msn.com +short || exit 1
    
CMD ["dnsmasq", "--no-daemon"]
