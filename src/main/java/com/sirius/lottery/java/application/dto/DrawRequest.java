package com.sirius.lottery.java.application.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class DrawRequest {
    private String instanceId;
    private String requestId;
}
