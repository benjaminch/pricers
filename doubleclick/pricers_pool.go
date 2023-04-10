package doubleclick

import "sync"

type PricersPool struct {
	pool                              *sync.Pool
	encryptionKeyRaw, integrityKeyRaw []byte
	scaleFactor                       float64
}

func NewPricersPool(encryptionKeyRaw, integrityKeyRaw []byte, scaleFactor float64) *PricersPool {
	return &PricersPool{
		pool:             &sync.Pool{},
		encryptionKeyRaw: encryptionKeyRaw,
		integrityKeyRaw:  integrityKeyRaw,
		scaleFactor:      scaleFactor,
	}
}

func (pl *PricersPool) AcquirePricer() *DoubleClickPricer {
	var pricer *DoubleClickPricer
	oldPricer := pl.pool.Get()
	if oldPricer != nil {
		pricer = oldPricer.(*DoubleClickPricer)
	} else {
		pricer = NewDoubleClickPricerFromRawKeys(pl.encryptionKeyRaw, pl.integrityKeyRaw, pl.scaleFactor)
	}
	return pricer
}

func (pl *PricersPool) ReleasePricer(pricer *DoubleClickPricer) {
	pl.pool.Put(pricer)
}
