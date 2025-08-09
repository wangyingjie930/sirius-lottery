package com.sirius.lottery.java.domain.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 中奖记录实体
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "win_record")
public class WinRecord {

    /**
     * 自增ID
     */
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    /**
     * 请求ID
     */
    @Column(name = "request_id")
    private String requestId;

    /**
     * 唯一订单号
     */
    @Column(name = "order_id", unique = true, nullable = false)
    private String orderId;

    /**
     * 抽奖实例ID
     */
    @Column(name = "instance_id", nullable = false)
    private String instanceId;

    /**
     * 用户ID
     */
    @Column(name = "user_id", nullable = false)
    private Long userId;

    /**
     * 奖品ID
     */
    @Column(name = "prize_id", nullable = false)
    private String prizeId;

    /**
     * 发放状态: 1-待发放, 2-发放成功, 3-发放失败
     */
    @Column(name = "status", nullable = false)
    private Integer status;
}
