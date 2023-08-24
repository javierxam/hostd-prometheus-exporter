package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	//	"fmt"
	"log"
    "time"
	"math/big"
    "go.sia.tech/hostd/api"
	"go.sia.tech/core/types"
	)

var (
/*	// Revenue Metrics
	hostStorageRevenue = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_storage_potential", Help: "Storage potential revenue"})
*/
	hostdTotalStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_total_storage", Help: "total amount of storage available on the hostd in bytes"})
	hostdUsedStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_used_storage", Help: "total amount of storage used on the hostd in bytes"})
	hostdRemainingStorage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_remaining_storage", Help: "amount of storage remaining on the host in bytes"})
	
	hostdIngress = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_ingress", Help: "Ingress potential revenue"})
	hostdEgress = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_egress", Help: "Egress potential revenue"})

	hostdLockedCollateral = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_locked_collateral", Help: "Locked collateral"})
	hostdRiskedCollateral = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "host_risked_collateral", Help: "Risked collateral"})

	walletConfirmedSiacoinBalance = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallet_confirmed_siacoin_balance", Help: "Wallet confirmed SCP balance"})

	hostdContractCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hostd_contract_count", Help: "number of host contracts"})
	)

//	Float64 converts types.Currency to float64.
	func Float64(c types.Currency) float64 {
    f, _ := new(big.Rat).SetInt(c.Big()).Float64()
    return f
	}

func callClient(passwd string, address string){
	client := api.NewClient("http://"+address+"/api", passwd)
    metrics, err := client.Metrics(time.Now())
    if err != nil {
        log.Fatalln(err)
    }

	hostdTotalStorage.Set(float64((metrics.Storage.TotalSectors)*4194304))
//	fmt.Println(metrics.Storage.TotalSectors*4194304)
	hostdUsedStorage.Set(float64((metrics.Storage.PhysicalSectors)*4194304))
//	fmt.Println((metrics.Storage.PhysicalSectors*4194304))
	hostdRemainingStorage.Set(float64((metrics.Storage.TotalSectors-metrics.Storage.PhysicalSectors)*4194304))
//	fmt.Println((metrics.Storage.TotalSectors-metrics.Storage.PhysicalSectors)*4194304)

	hostdIngress.Set(float64(metrics.Data.RHP2.Ingress+metrics.Data.RHP3.Ingress))
//	fmt.Println(metrics.Data.RHP2.Ingress+metrics.Data.RHP3.Ingress)
	hostdEgress.Set(float64(metrics.Data.RHP2.Egress+metrics.Data.RHP3.Egress))
//	fmt.Println(metrics.Data.RHP2.Egress+metrics.Data.RHP3.Egress)
	
	hostdLockedCollateral.Set((Float64(metrics.Contracts.LockedCollateral))/1e24)
	hostdRiskedCollateral.Set((Float64(metrics.Contracts.RiskedCollateral))/1e24)

	walletConfirmedSiacoinBalance.Set((Float64(metrics.Balance))/1e24)

	hostdContractCount.Set(float64(metrics.Contracts.Active))


}


