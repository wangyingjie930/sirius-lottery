package com.sirius.lottery.java.application.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class AssetRequest {
    private Integer cost;
    private Long userId;
}
