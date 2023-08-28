package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	//	"fmt"
	rhp2 "go.sia.tech/core/rhp/v2"
	"go.sia.tech/core/types"
	"go.sia.tech/hostd/api"
	"log"
	"math/big"
	"time"
)

var (
	/*	// Revenue Metrics
		hostStorageRevenue = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "host_storage_potential", Help: "Storage potential revenue"})
	*/
	hostdTotalStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_total_storage", Help: "Total amount of storage available on the hostd in bytes"})
	hostdUsedStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_used_storage", Help: "Total amount of storage used on the hostd in bytes"})
	hostdRemainingStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_remaining_storage", Help: "Amount of storage remaining on the host in bytes"})
	contractStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_contract_storage", Help: "Amount of contract storage on the host in bytes"})
	tempStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_temp_storage", Help: "Amount of temporary storage on the host in bytes"})

	storageReads = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_storage_reads", Help: "Amount of read operations"})
	storageWrites = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_storage_writes", Help: "Amount of write operations"})

	hostdIngress = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_ingress", Help: "Total ingress bandwidth usage"})
	hostdEgress = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_egress", Help: "Total egress bandwidth usage"})

	hostdLockedCollateral = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_locked_collateral", Help: "Locked collateral"})
	hostdRiskedCollateral = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_risked_collateral", Help: "Risked collateral"})

	walletConfirmedSiacoinBalance = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_wallet_confirmed_siacoin_balance", Help: "Wallet confirmed SCP balance"})

	hostdContractCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_contract_count", Help: "Number of host contracts"})
)

func convertCurrency(c types.Currency) float64 {
	f, _ := new(big.Rat).SetFrac(c.Big(), types.Siacoins(1).Big()).Float64()
	return f
}

func callClient(passwd string, address string) {
	client := api.NewClient("http://"+address+"/api", passwd)
	metrics, err := client.Metrics(time.Now())
	if err != nil {
		log.Fatalln(err)
	}

	hostdTotalStorage.Set(float64((metrics.Storage.TotalSectors) * rhp2.SectorSize))
	hostdUsedStorage.Set(float64((metrics.Storage.PhysicalSectors) * rhp2.SectorSize))
	contractStorage.Set(float64((metrics.Storage.ContractSectors) * rhp2.SectorSize))
	tempStorage.Set(float64((metrics.Storage.TempSectors) * rhp2.SectorSize))

	storageReads.Set(float64(metrics.Storage.Reads))
	storageWrites.Set(float64(metrics.Storage.Writes))

	hostdRemainingStorage.Set(float64((metrics.Storage.TotalSectors - metrics.Storage.PhysicalSectors) * rhp2.SectorSize))

	hostdIngress.Set(float64(metrics.Data.RHP2.Ingress + metrics.Data.RHP3.Ingress))
	hostdEgress.Set(float64(metrics.Data.RHP2.Egress + metrics.Data.RHP3.Egress))

	hostdLockedCollateral.Set(convertCurrency(metrics.Contracts.LockedCollateral))
	hostdRiskedCollateral.Set(convertCurrency(metrics.Contracts.RiskedCollateral))

	walletConfirmedSiacoinBalance.Set(convertCurrency(metrics.Balance))

	hostdContractCount.Set(float64(metrics.Contracts.Active))

}
