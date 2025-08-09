package com.sirius.lottery.java.application.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class CreateWinRecordRequest {
    private String orderId;
    private String prizeId;
    private String instanceId;
    private String requestId;
    private Long userId;
}
