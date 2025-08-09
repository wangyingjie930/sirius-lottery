package com.sirius.lottery.java.application.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class StockActionRequest {
    private String instanceId;
    private String prizeId;
    private Integer num;
}
