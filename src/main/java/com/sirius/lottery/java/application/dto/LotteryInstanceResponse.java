package com.sirius.lottery.java.application.dto;

import lombok.Data;

import java.time.LocalDateTime;
import java.util.List;

@Data
public class LotteryInstanceResponse {
    private String instanceId;
    private String name;
    private LocalDateTime startTime;
    private LocalDateTime endTime;
    private LocalDateTime serverTime;
    private String templateStyle;
    private TemplateConfig templateConfig;
    private List<Pool> pools;

    @Data
    public static class TemplateConfig {
        private String backgroundImage;
    }

    @Data
    public static class Pool {
        private String poolId;
        private String poolName;
        private List<Cost> cost;
        private List<Prize> prizes;
    }

    @Data
    public static class Cost {
        private String assetId;
        private Integer amount;
    }

    @Data
    public static class Prize {
        private Integer position;
        private String prizeId;
    }
}
