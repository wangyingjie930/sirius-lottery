package com.sirius.lottery.java.domain.model;

import com.sirius.lottery.java.domain.entity.Pool;
import com.sirius.lottery.java.domain.entity.Prize;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class DrawContext {
    private String instanceId;
    private Long userId;
    private Pool pool;
    private List<Prize> prizes;
}
