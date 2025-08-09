package com.sirius.lottery.java.application.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class DrawResponse {
    private String orderId;
    private String prizeId;
    private Boolean isWin;
}
