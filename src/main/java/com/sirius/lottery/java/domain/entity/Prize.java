package com.sirius.lottery.java.domain.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 奖品实体
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "prize")
public class Prize {

    /**
     * 自增ID
     */
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    /**
     * 奖池ID
     */
    @Column(name = "pool_id")
    private Long poolId;

    /**
     * 业务奖品ID, 来自奖品中心
     */
    @Column(name = "prize_id")
    private String prizeId;

    /**
     * 奖品名称
     */
    @Column(name = "prize_name")
    private String prizeName;

    /**
     * 总预算库存
     */
    @Column(name = "allocated_stock")
    private Integer allocatedStock;

    /**
     * 概率
     */
    @Column(name = "probability")
    private Double probability;

    /**
     * 是否特殊奖品(如谢谢参与)
     */
    @Column(name = "is_special")
    private Boolean isSpecial;
}
