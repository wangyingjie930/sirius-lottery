package com.sirius.lottery.java.interfaces.web;

import com.sirius.lottery.java.application.dto.DrawRequest;
import com.sirius.lottery.java.application.dto.DrawResponse;
import com.sirius.lottery.java.application.dto.LotteryInstanceResponse;
import com.sirius.lottery.java.application.service.LotteryService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v2/lottery")
@RequiredArgsConstructor
public class LotteryController {

    private final LotteryService lotteryService;

    @PostMapping("/draw")
    public ResponseEntity<DrawResponse> draw(@RequestBody DrawRequest request) {
        DrawResponse response = lotteryService.draw(request);
        return ResponseEntity.ok(response);
    }

    @GetMapping("/instance/{instanceId}")
    public ResponseEntity<LotteryInstanceResponse> getLotteryInstance(@PathVariable String instanceId) {
        LotteryInstanceResponse response = lotteryService.getLotteryInstance(instanceId);
        return ResponseEntity.ok(response);
    }
}
