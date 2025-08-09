package com.sirius.lottery.java.domain.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;
import java.util.List;

/**
 * 抽奖实例实体
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "instance")
public class Instance {

    /**
     * 自增ID
     */
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    /**
     * 业务活动ID
     */
    @Column(name = "instance_id", unique = true, nullable = false)
    private String instanceId;

    /**
     * 活动名称
     */
    @Column(name = "instance_name", nullable = false)
    private String instanceName;

    /**
     * 关联的模板ID
     */
    @Column(name = "template_id", nullable = false)
    private Long templateId;

    /**
     * 活动开始时间
     */
    @Column(name = "start_time", nullable = false)
    private LocalDateTime startTime;

    /**
     * 活动结束时间
     */
    @Column(name = "end_time", nullable = false)
    private LocalDateTime endTime;

    /**
     * 状态: 1-待上线, 2-进行中, 3-已下线
     */
    @Column(name = "status", nullable = false)
    private Integer status;

    /**
     * 奖池列表
     */
    @Transient
    private List<Pool> pools;
}
