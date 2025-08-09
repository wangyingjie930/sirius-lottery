package com.sirius.lottery.java.domain.strategy;

import com.sirius.lottery.java.domain.entity.Prize;
import com.sirius.lottery.java.domain.model.DrawContext;

public interface LotteryStrategy {
    Prize draw(DrawContext drawContext);
}
