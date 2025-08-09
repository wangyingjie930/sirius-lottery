package com.sirius.lottery.java.domain.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

/**
 * 奖池实体
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "pool")
public class Pool {

    /**
     * 自增ID
     */
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    /**
     * 抽奖实例ID
     */
    @Column(name = "instance_id")
    private String instanceId;

    /**
     * 奖池名称
     */
    @Column(name = "pool_name")
    private String poolName;

    /**
     * 消耗的资产列表, e.g., [{"asset_id": "ticket", "amount": 1}]
     */
    @Column(name = "cost_json", columnDefinition = "json")
    private String costJson;

    /**
     * 抽奖算法策略
     */
    @Column(name = "lottery_strategy")
    private String lotteryStrategy;

    /**
     * 策略相关配置, e.g., {"guarantee_count": 10}
     */
    @Column(name = "strategy_config_json", columnDefinition = "json")
    private String strategyConfigJson;

    /**
     * 奖品列表
     */
    @Transient
    private List<Prize> prizes;
}
